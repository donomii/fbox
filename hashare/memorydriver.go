package hashconnect

import (
	"encoding/hex"
	"log"
	"regexp"
	"strings"

	"bytes"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/donomii/fbox"
	"github.com/donomii/hashare"
)

type HashareDriver struct {
	Store     *hashare.SqlStore
	BlockSize int
	Files     map[string]*HashareFile
	Username  string
	Password  string
}

func (d *HashareDriver) Authenticate(username string, password string) bool {
	return true
	return username == d.Username && password == d.Password
}

func (d *HashareDriver) Bytes(path string) int64 {
	if f, ok := d.Files[path]; ok {
		return f.File.Size()
	} else {
		return -1
	}
}

func (d *HashareDriver) ModifiedTime(path string) (time.Time, bool) {
	if f, ok := d.Files[path]; ok {
		return f.File.ModTime(), true
	} else {
		return time.Now(), false
	}
}

func (d *HashareDriver) ChangeDir(path string) bool {

	return true
}

func (d *HashareDriver) DirContents(path string) ([]os.FileInfo, bool) {
	log.Println("Fetching directory contents for", path)
	path = strings.TrimSuffix(path, "/-l")
	pathlets, ok := hashare.ResolvePath(d.Store, []byte(path), d.BlockSize)
	if !ok {
		log.Println("Could not find directory, returning error")
		return nil, false
	}
	log.Println("Pathlets:", hashare.BytesArrayToString(pathlets))

	//Get the name of our current working directory
	currentDir := pathlets[len(pathlets)-1]

	files := []os.FileInfo{}
	dir := hashare.FetchDirectory(d.Store, currentDir, d.BlockSize)
	for i, v := range dir.Entries {
		log.Printf("%v: %v (%v)\n", i, string(v.Name), hex.Dump(v.Id))

		if string(v.Type) == "dir" {
			f := fbox.NewDirItem(string(v.Name))
			files = append(files, f)
		} else {
			f := fbox.NewFileItem(string(v.Name), 10, time.Now().UTC())
			files = append(files, f)
		}
	}
	//Windows freaks out if it doesn't get at list one file in the file list
	if len(files) == 0 {
		f := fbox.NewDirItem(".")
		files = append(files, f)
		f = fbox.NewDirItem("..")
		files = append(files, f)
	}
	sort.Sort(&FilesSorter{files})
	return files, true
}

func (d *HashareDriver) DeleteDir(path string) bool {
	log.Println("Deleting directory", path)
	pathlets, ok := hashare.ResolvePath(d.Store, []byte(path), d.BlockSize)
	if !ok {
		//Deleting a non-existant directory still counts as a win
		return true
	}
	//log.Println("Pathlets:", hashare.BytesArrayToString(pathlets))

	//Hashare treats files and directories mostly the same
	hashare.DeleteFile(d.Store, pathlets, d.BlockSize, true)
	return true
}

func (d *HashareDriver) DeleteFile(path string) bool {
	log.Println("Deleting file", path)
	pathlets, ok := hashare.ResolvePath(d.Store, []byte(path), d.BlockSize)
	if !ok {
		//Deleting a file that doesn't exist still counts imo
		return true
	}
	log.Println("Pathlets:", hashare.BytesArrayToString(pathlets))
	hashare.DeleteFile(d.Store, pathlets, d.BlockSize, true)
	return true
}

func (d *HashareDriver) Rename(from_path string, to_path string) bool {
	from_pathlets := regexp.MustCompile("\\\\|/").Split(from_path, -1)
	to_pathlets := regexp.MustCompile("\\\\|/").Split(to_path, -1)
	filename := to_pathlets[len(to_pathlets)-1]
	to_pathlets = to_pathlets[0 : len(to_pathlets)-1]
	hashare.MoveFile(d.Store, filename, hashare.StringArrayToBytes(from_pathlets), hashare.StringArrayToBytes(to_pathlets), d.BlockSize, true)
	return true
}

func (d *HashareDriver) MakeDir(path string) bool {
	log.Println("Making directory", path)
	pathlets, ok := hashare.ResolvePath(d.Store, []byte(path), d.BlockSize)
	if ok {
		log.Println("Directory already exists!")
		//Creating a directory that already exists still counts
		//We should probably update the mtime or something?
		return true
	}
	splits := regexp.MustCompile("\\\\|/").Split(path, -1)
	filename := splits[len(splits)-1]
	//pathlets = pathlets[0:len(pathlets)-1]
	hashare.MkDir(d.Store, pathlets, filename, d.BlockSize)
	return true
}

func (d *HashareDriver) GetFile(path string, position int64) (io.ReadCloser, bool) {
	log.Println("Fetching file", path)
	pathlets, ok := hashare.ResolvePath(d.Store, []byte(path), d.BlockSize)
	if !ok {
		return nil, ok
	}
	log.Println("Pathlets:", hashare.BytesArrayToString(pathlets))

	//Get the name of our current working directory
	currentDir := pathlets[len(pathlets)-1]

	raw_content := hashare.FetchFile(d.Store, currentDir, d.BlockSize)

	return ioutil.NopCloser(bytes.NewReader(raw_content[position:])), true
}

func (d *HashareDriver) PutFile(path string, reader io.Reader) bool {
	log.Println("fbox: Putting file", path)
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("fbox: Error reading file data:", err)
		return false
	}

	pathlets, ok := hashare.ResolvePath(d.Store, []byte(path), d.BlockSize)
	if ok {
		log.Println("fbox: File already exists!")
		return !ok
	}
	//Get the name of our current working directory
	//currentDir := pathlets[len(pathlets)-1]

	splits := regexp.MustCompile("\\\\|/").Split(path, -1)
	filename := splits[len(splits)-1]

	//pathlets = pathlets[0:len(pathlets)-1]

	log.Println("fbox: Pathlets for putbytes:", hashare.BytesArrayToString(pathlets))
	hashare.PutBytes(d.Store, bytes, filename, pathlets, d.BlockSize, true)
	//d.Files[path] = &HashareFile{fbox.NewFileItem(filepath.Base(path), int64(len(bytes)), time.Now().UTC()), bytes}

	log.Println("fbox: Put file complete:", path)
	return true
}
