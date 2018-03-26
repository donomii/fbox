package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/donomii/hashare"
	"github.com/donomii/vort"
	"github.com/donomii/vort/hashare"
)

func main() {
	log.Println("Application:", os.Args[0])
	host := "127.0.0.1"
	port := 8021
	username := "test"
	password := "test"
	debug := true
	conf := hashare.Config{
		Store:          nil,
		Blocksize:      500,
		UseEncryption:  true,
		UseCompression: false,
		EncryptionKey:  []byte("a very very very very secret key"),
	}

	flag.IntVar(&conf.Blocksize, "blocksize", 1048576, "Store data in chunks of this size")
	flag.BoolVar(&conf.UseEncryption, "encrypt", false, "Encrypt every block")
	flag.BoolVar(&conf.UseCompression, "compress", true, "Compress every block")
	var optStr string
	var optStoreType string
	repository := os.Getenv("USERPROFILE") + "/Desktop/default.vort"
	flag.StringVar(&optStr, "key", "a very very very very secret key", "Encryption key")
	flag.StringVar(&optStoreType, "type", "auto", "Repository type (sql or files)")
	flag.StringVar(&repository, "repo", repository, "Path to repository directory")
	flag.Parse()
	conf.EncryptionKey = []byte(optStr)

	if flag.Arg(0) != "" {
		repository = flag.Arg(0)
	}
	files := map[string]*hashconnect.HashareFile{
		"/": &hashconnect.HashareFile{vort.NewDirItem("", 0, time.Now().UTC()), nil},
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
		} else {
			s = hashare.NewSQLStore(repository)
		}
	} else {

		if optStoreType == "files" {
			s = hashare.NewSQLStore(repository)
		} else {
			s = hashare.NewFileStore(repository)
		}
	}
	conf = s.Init(conf)
	log.Println("Opened repository:", repository)
	conf.Store = s
	log.Printf("Config: %+v", conf)
	factory := &hashconnect.HashareDriverFactory{conf, files, username, password}

	for {
		port = port + 1
		server := vort.NewFTPServer(&vort.FTPServerOpts{
			ServerName: "Example FTP server",
			Factory:    factory,
			Hostname:   host,
			Port:       port,
			PassiveOpts: &vort.PassiveOpts{
				ListenAddress: host,
				NatAddress:    host,
				PassivePorts: &vort.PassivePorts{
					Low:  42000,
					High: 45000,
				},
			},
		})

		go func() {
			time.Sleep(1 * time.Second)
			log.Printf("vort FTP server listening on %s:%d", host, port)
			log.Printf("Access: ftp://%s:%s@%s:%d/", username, password, host, port)

			cmd := exec.Command("c:/Windows/explorer.exe", fmt.Sprintf("ftp://%s:%s@%s:%d/", username, password, host, port))
			cmd.Start()
			log.Println("Launched explorer window")
		}()
		err := server.ListenAndServe()

		if err != nil {
			log.Println(err)
		}
	}
}
