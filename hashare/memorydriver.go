package hashconnect

import (
"regexp"
"log"
"fmt"
"encoding/hex"

	"bytes"
	"github.com/donomii/fbox"
	"github.com/donomii/hashare"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type HashareDriver struct {
	Store	*hashare.SqlStore
	BlockSize	int
	Files    map[string]*HashareFile
	Username string
	Password string
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
	if f, ok := d.Files[path]; ok && f.File.IsDir() {
		return true
	} else {
		return false
	}
}

func (d *HashareDriver) DirContents(path string) ([]os.FileInfo, bool) {
pathlets := hashare.ResolvePath(d.Store, []byte(path), d.BlockSize) 
	
	log.Println("Pathlets:", hashare.BytesArrayToString(pathlets))
	
	//Get the name of our current working directory
	currentDir := pathlets[len(pathlets)-1]

	files := []os.FileInfo{}
		dir := hashare.FetchDirectory(d.Store, currentDir, d.BlockSize)
		for i, v := range dir.Entries {
				fmt.Printf("%v: %v (%v)\n", i,string(v.Name), hex.Dump(v.Id))
				
				
				f:= fbox.NewFileItem(string(v.Name), 10, time.Now().UTC())
				files = append(files, f)
		}

		sort.Sort(&FilesSorter{files})

		return files, true
}

func (d *HashareDriver) DeleteDir(path string) bool {
	if f, ok := d.Files[path]; ok && f.File.IsDir() {
		haschildren := false
		for p, _ := range d.Files {
			if strings.HasPrefix(p, path+"/") {
				haschildren = true
				break
			}
		}

		if haschildren {
			return false
		}

		delete(d.Files, path)

		return true
	} else {
		return false
	}
}

func (d *HashareDriver) DeleteFile(path string) bool {
	if f, ok := d.Files[path]; ok && !f.File.IsDir() {
		delete(d.Files, path)
		return true
	} else {
		return false
	}
}

func (d *HashareDriver) Rename(from_path string, to_path string) bool {
	if f, from_path_exists := d.Files[from_path]; from_path_exists {
		if _, to_path_exists := d.Files[to_path]; !to_path_exists {
			if _, to_path_parent_exists := d.Files[filepath.Dir(to_path)]; to_path_parent_exists {
				if f.File.IsDir() {
					delete(d.Files, from_path)
					d.Files[to_path] = &HashareFile{fbox.NewDirItem(filepath.Base(to_path)), nil}
					torename := make([]string, 0)
					for p, _ := range d.Files {
						if strings.HasPrefix(p, from_path+"/") {
							torename = append(torename, p)
						}
					}
					for _, p := range torename {
						sf := d.Files[p]
						delete(d.Files, p)
						np := to_path + p[len(from_path):]
						d.Files[np] = sf
					}
				} else {
					delete(d.Files, from_path)
					d.Files[to_path] = &HashareFile{fbox.NewFileItem(filepath.Base(to_path), f.File.Size(), f.File.ModTime()), f.Content}
				}
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else {
		return false
	}
}

func (d *HashareDriver) MakeDir(path string) bool {
	if _, ok := d.Files[path]; ok {
		return false
	} else {
		d.Files[path] = &HashareFile{fbox.NewDirItem(filepath.Base(path)), nil}
		return true
	}
}

func (d *HashareDriver) GetFile(path string, position int64) (io.ReadCloser, bool) {
pathlets := hashare.ResolvePath(d.Store, []byte(path), d.BlockSize) 
	
	log.Println("Pathlets:", hashare.BytesArrayToString(pathlets))
	
	//Get the name of our current working directory
	currentDir := pathlets[len(pathlets)-1]

	raw_content := hashare.FetchFile(d.Store, currentDir, d.BlockSize)

	return ioutil.NopCloser(bytes.NewReader(raw_content[position:])), true
}

func (d *HashareDriver) PutFile(path string, reader io.Reader) bool {
	bytes, err := ioutil.ReadAll(reader)
			if err != nil {
				return false
			}

			pathlets := hashare.ResolvePath(d.Store, []byte(path), d.BlockSize) 
	
			log.Println("Pathlets:", hashare.BytesArrayToString(pathlets))
			
			//Get the name of our current working directory
			//currentDir := pathlets[len(pathlets)-1]
	
			//pathlets = pathlets[0:len(pathlets)-1]
			
			splits := regexp.MustCompile("\\\\|/").Split(path, -1)
			filename := splits[len(splits)-1]
			
			

			
	
			hashare.PutBytes(d.Store, bytes, filename, pathlets, d.BlockSize, true)
			//d.Files[path] = &HashareFile{fbox.NewFileItem(filepath.Base(path), int64(len(bytes)), time.Now().UTC()), bytes}

			return true
}
