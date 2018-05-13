package hashconnect

import (

	//"bytes"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	_ "strings"
	"time"

	"github.com/donomii/hashare"
	"github.com/donomii/vort-ftprelay"
)

type HashareDriver struct {
	Conf     *hashare.Config
	Files    map[string]*HashareFile
	Username string
	Password string
}

func (d *HashareDriver) Authenticate(username string, password string) bool {
	return username == d.Username && password == d.Password
}

func (d *HashareDriver) Bytes(path string) int64 {
	log.Println("vort: Fetching file size for", path)
	var ok bool
	var meta *hashare.DirectoryEntry
	hashare.WithTransaction(d.Conf, func(tr hashare.Transaction) hashare.Transaction {
		//Hashare treats files and directories mostly the same
		meta, ok = hashare.GetMeta(path, d.Conf, tr)
		return tr
	})

	if !ok {
		log.Println("vort: File not found, returning -1")
		return -1
	}
	log.Println("vort: Returning file size", meta.Size)
	if meta.Size < 0 {
		return 0
	}
	return meta.Size

}

func getCurrentMeta(path string, conf *hashare.Config) (*hashare.DirectoryEntry, bool) {
	var ok bool
	var f *hashare.DirectoryEntry
	hashare.WithTransaction(conf, func(tr hashare.Transaction) hashare.Transaction {
		//Hashare treats files and directories mostly the same
		f, ok = hashare.GetMeta(path, conf, tr)
		return tr
	})
	return f, ok
}
func (d *HashareDriver) ModifiedTime(path string) (time.Time, bool) {

	f, ok := getCurrentMeta(path, d.Conf)
	if ok {
		return f.Modified, true
	} else {
		t1, _ := time.Parse(time.RFC3339, "1981-11-01T22:08:41+00:00")
		return t1, false
	}
}

func (d *HashareDriver) ChangeDir(path string) bool {
	//Maybe we should have a DirectoryExists() API?
	log.Println("Changing to path", path)

	var ok bool

	hashare.WithTransaction(d.Conf, func(tr hashare.Transaction) hashare.Transaction {
		//Hashare treats files and directories mostly the same
		_, ok = hashare.GetMeta(path, d.Conf, tr)
		return tr
	})
	return ok
}

func (d *HashareDriver) DirContents(path string) ([]os.FileInfo, bool) {
	log.Println("Fetching directory contents for", path)
	//We have to return a list of files in the operating system format
	files := []os.FileInfo{}

	dirEntries, ok := hashare.List(d.Conf.Store, path, d.Conf)
	if !ok {
		log.Println("Could not find directory, returning error")
		return nil, false
	}
	for i, v := range dirEntries {
		log.Printf("%v: %v (%v)\n", i, string(v.Name), hex.Dump(v.Id))

		if string(v.Type) == "dir" {
			f := vort.NewDirItem(string(v.Name), v.Size, time.Now().UTC())
			files = append(files, f)
		} else {
			f := vort.NewFileItem(string(v.Name), v.Size, time.Now().UTC())
			files = append(files, f)
		}
	}
	//Windows freaks out if it doesn't get at least one file in the file list
	if len(files) == 0 {
		f := vort.NewDirItem(".", 0, time.Now().UTC())
		files = append(files, f)
		f = vort.NewDirItem("..", 0, time.Now().UTC())
		files = append(files, f)
	}
	sort.Sort(&FilesSorter{files})
	return files, true
}

func (d *HashareDriver) DeleteDir(path string) bool {
	log.Println("vort: Deleting directory:", path)

	//FTP thinks that it is an error to delete a file that doesn't exist
	//Mission accomplished IMO, but we must follow the spec...
	meta, ok := getCurrentMeta(path, d.Conf)
	if !ok {
		return false
	}
	dirEntries, ok := hashare.List(d.Conf.Store, path, d.Conf)
	if len(dirEntries) > 0 {
		//Can't delete a directory while it still contains items, apparently
		return false
	}
	if string(meta.Type) != "dir" {
		//We need different commands to delete files and directories, because reasons
		return false
	}
	hashare.WithTransaction(d.Conf, func(tr hashare.Transaction) hashare.Transaction {
		//Hashare treats files and directories mostly the same
		ret, _ := hashare.DeleteFile(d.Conf.Store, path, d.Conf, false, tr)
		return ret
	})
	return true
}

func (d *HashareDriver) DeleteFile(path string) bool {
	log.Println("vort: Deleting file:", path)
	//FTP thinks that it is an error to delete a file that doesn't exist
	//Mission accomplished IMO, but we must follow the spec...
	meta, ok := getCurrentMeta(path, d.Conf)
	if !ok {
		return false
	}
	if string(meta.Type) != "file" {
		//We need different commands to delete files and directories, because reasons
		return false
	}
	hashare.WithTransaction(d.Conf, func(tr hashare.Transaction) hashare.Transaction {
		ret, _ := hashare.DeleteFile(d.Conf.Store, path, d.Conf, false, tr)
		return ret
	})
	return true
}

func (d *HashareDriver) Rename(from_path string, to_path string) bool {
	//This also should be successful, but the spec says otherwise....
	if from_path == to_path {
		return false
	}
	ok := hashare.MoveFile(d.Conf.Store, from_path, to_path, d.Conf, true)
	return ok
}

func (d *HashareDriver) MakeDir(path string) bool {
	log.Println("Making directory", path)
	_, ok := getCurrentMeta(path, d.Conf)
	if ok {
		//This also should be success, but the spec...
		return false
	}
	//pathlets = pathlets[0:len(pathlets)-1]
	hashare.MkDir(d.Conf.Store, path, d.Conf)
	return true
}

func (d *HashareDriver) GetFile(path string, position int64) (io.ReadCloser, bool) {
	reader, ok := hashare.GetFileStream(d.Conf.Store, path, position, -1, d.Conf)
	//FIXME reader is actually a readcloser, we just have the type wrong on GetFileStream
	return ioutil.NopCloser(reader), ok
}

func (d *HashareDriver) PutFile(path string, reader io.Reader) bool {
	log.Println("vort: Putting file", path)

	//pathlets = pathlets[0:len(pathlets)-1]

	//log.Println("vort: Pathlets for putbytes:", hashare.BytesArrayToString(pathlets))
	var ok bool
	hashare.WithTransaction(d.Conf, func(tr hashare.Transaction) (ret hashare.Transaction) {
		ret, ok = hashare.PutStream(d.Conf.Store, reader, path, d.Conf, true, tr)

		return
	})
	//d.Files[path] = &HashareFile{fbox.NewFileItem(filepath.Base(path), int64(len(bytes)), time.Now().UTC()), bytes}
	if ok {
		log.Println("vort: Put file complete:", path)
	} else {
		log.Println("vort: Put file failed:", path)
	}
	return ok
}
