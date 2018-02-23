package main

import (
"io/ioutil"
	"github.com/donomii/fbox"
	"github.com/donomii/fbox/hashare"
	"github.com/donomii/hashare"
	"log"
)

func main() {
	host := "127.0.0.1"
	port := 8021
	username := "test"
	password := "test"
	debug := false
	
	files := map[string]*hashconnect.HashareFile{
		"/": &hashconnect.HashareFile{fbox.NewDirItem(""), nil},
	}
	
		//Switch log output off by default
	if !debug {
		log.SetOutput(ioutil.Discard)
		log.SetFlags(0)
	}
	
	repository := "repository.cas"
	//Open the repository
	s := hashare.NewSQLStore(repository)
	s.Init()
	log.Println("Opened repository")
	
	

	factory := &hashconnect.HashareDriverFactory{s, 10000000, files, username, password}

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

	log.Printf("Example FTP server listening on %s:%d", host, port)
	log.Printf("Access: ftp://%s:%s@%s:%d/", username, password, host, port)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
