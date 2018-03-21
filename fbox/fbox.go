package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/donomii/fbox"
	"github.com/donomii/fbox/hashare"
	"github.com/donomii/hashare"
)

func main() {
	host := "127.0.0.1"
	port := 8021
	username := "test"
	password := "test"
	debug := true
	conf := hashare.Config{
		Store:          nil,
		Blocksize:      499,
		UseEncryption:  true,
		UseCompression: false,
		EncryptionKey:  []byte("a very very very very secret key"),
	}

	flag.IntVar(&conf.Blocksize, "blocksize", 500, "Store data in chunks of this size")
	flag.BoolVar(&conf.UseEncryption, "encrypt", false, "Encrypt every block")
	flag.BoolVar(&conf.UseCompression, "compress", true, "Compress every block")
	var optStr string
	var optStoreType string
	repository := "repository.cas"
	flag.StringVar(&optStr, "key", "a very very very very secret key", "Encryption key")
	flag.StringVar(&optStoreType, "type", "auto", "Repository type (sql or files)")
	flag.StringVar(&repository, "repo", "filebox.fbox", "Path to repository directory")
	flag.Parse()
	conf.EncryptionKey = []byte(optStr)

	if flag.Arg(0) != "" {
		repository = flag.Arg(0)
	}
	files := map[string]*hashconnect.HashareFile{
		"/": &hashconnect.HashareFile{fbox.NewDirItem("", 0, time.Now().UTC()), nil},
	}

	//Switch log output off by default
	if !debug {
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)
	}
	var s hashare.SiloStore
	//Open the repository
	if optStoreType == "auto" {
		//If the file exists, autodetect and open it
		if stat, err := os.Stat(repository); err == nil {
			if stat.Mode().IsDir() {
				//It's a fileblocks repo
				s = hashare.NewFileStore(repository)
			} else {
				//It's an SQLite filestore
				s = hashare.NewSQLStore(repository)
			}
		}
	} else {

		if optStoreType == "sql" {
			s = hashare.NewSQLStore(repository)
		} else {
			s = hashare.NewFileStore(repository)
		}
	}
	s.Init()
	log.Println("Opened repository")
	conf.Store = s
	log.Printf("Config: %+v", conf)
	factory := &hashconnect.HashareDriverFactory{conf, files, username, password}

	for {
		port = port + 1
		server := fbox.NewFTPServer(&fbox.FTPServerOpts{
			ServerName: "Example FTP server",
			Factory:    factory,
			Hostname:   host,
			Port:       port,
			PassiveOpts: &fbox.PassiveOpts{
				ListenAddress: host,
				NatAddress:    host,
				PassivePorts: &fbox.PassivePorts{
					Low:  42000,
					High: 45000,
				},
			},
		})

		go func() {
			log.Printf("FBOX FTP server listening on %s:%d", host, port)
			log.Printf("Access: ftp://%s:%s@%s:%d/", username, password, host, port)

			cmd := exec.Command("explorer.exe", fmt.Sprintf("ftp://%s:%s@%s:%d/", username, password, host, port))
			cmd.Start()
		}()
		err := server.ListenAndServe()

		if err != nil {
			log.Println(err)
		}
	}
}
