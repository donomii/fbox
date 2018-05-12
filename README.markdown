[![Build Status](https://travis-ci.org/donomii/vort-ftprelay.svg?branch=master)](https://travis-ci.org/donomii/vort-ftprelay)
[![GoDoc](https://godoc.org/github.com/donomii/vort-ftprelay?status.svg)](https://godoc.org/github.com/donomii/vort-ftprelay)

# Vort-ftprelay

This is part of the [Vort](http://github.com/donomii/vort) project.

vort-ftprelay provides access to vort files and vort fileservers through a FTP server.  Vort-relay starts an ftp server, that you connect to.  Your requests for files are translated to the native vort protocol.

You will get much better performance from the command line client, or the native filesystem mount.  This relay is provided for situations where neither of those will work for you.

## Installation

### Windows

Download [Vort](https://github.com/donomii/vort/releases/latest)

### Linux and MacOSX

Install google's go language, then:

go get github.com/donomii/Vort
go build github.com/donomii/Vort/vort/vort.go


## Command Line Examples

Start Vort with the default options

    vort

Connect to a network share

    	vort.exe --type=http --repo=http://localhost:80/
		
Open a vort file

	vort myfiles.vort
        
Start Vort with encryption

    vort --encrypt=1 --key="a 32-byte key123a 32-byte key123"
    
 Connect to Vort with an ftp client
 
     ftp localhost 8021
	
Create and open a sqlite type store

	vort --type=sqlite myfiles.vort

Create and open a file blocks type store

	vort --type=files myfiles.vort	
	

## Accessing Vort

Vort provides access to your vort file store through a local FTP server. Vort will print the url of the server, after it starts up.

You can access this url through normal FTP clients, including:

## Web browsers in FTP mode

Most web browsers include an ftp client. If you have Microsoft Edge, Firefox or Chrome, just paste the vort url into the address bar to access your files.

## FTP "drive" mount 

Windows has built in support for using FTP servers as ordinary drives. In a file explorer window, click on your computer, then find "Map network drive" somewhere in the menu. Follow the instructions to connect to your vort file.

Linux can also mount FTP drives, but it requires installing software and some fiddling about on the command line.

## Stand alone FTP clients

There are a large selection of FTP clients for every platform, and they should all work with Vort.

# Store types

Vort supports many different formats, which allows me to test many different storage techniques.  There are currently three useful ones

## SQLite

vort can use a single-file sqlite database as a file store.  This is the best option if you want to send your vort store to someone else, upload it to a server, or otherwise move it around.

It tends to slow down after you put ~100Gb into it, so if you want to store something bigger, try the files store.

	vort --type=sqlite myfiles.vort

Note: you only have to specify the type once, when creating your vort file.  vort will autodetect the store type of an existing file.

## Files

vort keeps its data in ordinary files.  This is probably slower than sqlite, but can handle almost any amount of data without slowing down (as much).  This is normally the best option for network servers (it is the default for vort-nfs).

	vort --type=files myfiles.vort

Note: you only have to specify the type once, when creating your vort directory.  vort will autodetect the type of an existing file.

## http

Technically not a file store, the http driver connects to a vort network share, and accesses the data on the server.

	vort --type=http --repo=http://myserver.localnet/

The slash on the end is required, for now.

# Encryption

Warning: Vort uses strong encryption. If you lose your encryption key, you will never be able to access your files again.

Vort uses symmetric key encryption, so in order to access your encrypted filesystem, you must remember the key you created, and type it in every time you open Vort. If you forget or lose your key, all your files are gone forever. They are never coming back. Using encryption is entirely at your own risk, and I am not responsible in any way for any data loss or leak due to using encryption.

# Disclaimer

Actually, I'm not responsible at all for any problems caused by using this software. It is open source, you're getting it for free, and if it doesn't work you're entitled to a full refund and nothing more.
