package hashconnect

import (
	"encoding/hex"
	"log"
	"regexp"
	_ "strings"

	"bytes"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/donomii/hashare"
	"github.com/donomii/vort"
)

type HashareDriver struct {
	Conf     *hashare.Config
	Files    map[string]*HashareFile
	Username string
	Password string
}

func (d *HashareDriver) Authenticate(username string, password string) bool {
	return true //FIXME
	return username == d.Username && password == d.Password
}

func (d *HashareDriver) Bytes(path string) int64 {
	log.Println("vort: Fetching file size for", path)
	meta, ok := hashare.GetMeta(d.Conf.Store, path, d.Conf)
	if !ok {
		log.Println("vort: File not found, returning -1")
		return -1
	}
	log.Println("vort: Returning file size", meta.Size)
	return meta.Size
}

func (d *HashareDriver) ModifiedTime(path string) (time.Time, bool) {
	if f, ok := d.Files[path]; ok {
		return f.File.ModTime(), true
	} else {
		t1, _ := time.Parse(time.RFC3339, "1981-11-01T22:08:41+00:00")
		return t1, true
	}
}

func (d *HashareDriver) ChangeDir(path string) bool {
	_, ok := hashare.ResolvePath(d.Conf.Store, []byte(path), d.Conf)
	return ok
}

func (d *HashareDriver) DirContents(path string) ([]os.FileInfo, bool) {
	log.Println("Fetching directory contents for", path)
	//path = strings.TrimSuffix(path, "/-l")
	//path = strings.TrimSuffix(path, "/ls-l")
	pathlets, ok := hashare.ResolvePath(d.Conf.Store, []byte(path), d.Conf)
	if !ok {
		log.Println("Could not find directory, returning error")
		return nil, false
	}
	log.Println("Pathlets:", hashare.BytesArrayToString(pathlets))

	//Get the name of our current working directory
	currentDir := pathlets[len(pathlets)-1]

	files := []os.FileInfo{}
	dir := hashare.FetchDirectory(d.Conf.Store, currentDir, d.Conf)
	for i, v := range dir.Entries {
		log.Printf("%v: %v (%v)\n", i, string(v.Name), hex.Dump(v.Id))

		if string(v.Type) == "dir" {
			f := vort.NewDirItem(string(v.Name), v.Size, time.Now().UTC())
			files = append(files, f)
		} else {
			f := vort.NewFileItem(string(v.Name), v.Size, time.Now().UTC())
			files = append(files, f)
		}
	}
	//Windows freaks out if it doesn't get at list one file in the file list
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
	//Hashare treats files and directories mostly the same
	hashare.DeleteFile(d.Conf.Store, path, d.Conf, true)
	return true
}

func (d *HashareDriver) DeleteFile(path string) bool {
	log.Println("vort: Deleting file:", path)
	hashare.DeleteFile(d.Conf.Store, path, d.Conf, true)
	return true
}

func (d *HashareDriver) Rename(from_path string, to_path string) bool {

	hashare.MoveFile(d.Conf.Store, from_path, to_path, d.Conf, true)
	return true
}

func (d *HashareDriver) MakeDir(path string) bool {
	log.Println("Making directory", path)

	//pathlets = pathlets[0:len(pathlets)-1]
	hashare.MkDir(d.Conf.Store, path, d.Conf)
	return true
}

func (d *HashareDriver) GetFile(path string, position int64) (io.ReadCloser, bool) {
	content, ok := hashare.GetFile(d.Conf.Store, path, position, d.Conf)
	return ioutil.NopCloser(bytes.NewReader(content)), ok
}

func (d *HashareDriver) PutFile(path string, reader io.Reader) bool {
	log.Println("vort: Putting file", path)

	pathlets, ok := hashare.ResolvePath(d.Conf.Store, []byte(path), d.Conf)
	if ok {
		log.Println("vort: File already exists!")
		return !ok
	}
	//Get the name of our current working directory
	//currentDir := pathlets[len(pathlets)-1]

	splits := regexp.MustCompile("\\\\|/").Split(path, -1)
	filename := splits[len(splits)-1]

	//pathlets = pathlets[0:len(pathlets)-1]

	//log.Println("vort: Pathlets for putbytes:", hashare.BytesArrayToString(pathlets))
	hashare.PutStream(d.Conf.Store, reader, filename, pathlets, d.Conf, true)
	//d.Files[path] = &HashareFile{fbox.NewFileItem(filepath.Base(path), int64(len(bytes)), time.Now().UTC()), bytes}

	log.Println("vort: Put file complete:", path)
	return true
}
