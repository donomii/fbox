[![Build Status](https://travis-ci.org/donomii/Vort.svg?branch=master)](https://travis-ci.org/donomii/Vort)
[![GoDoc](https://godoc.org/github.com/donomii/Vort?status.svg)](https://godoc.org/github.com/donomii/Vort)

# Vort

A secure, encrypted file store.

## Download

Warning:  this is experimental software.  Do not use it for anything important.  Make backups often!

[Windows](https://github.com/donomii/Vort/releases)

[Linux](https://github.com/donomii/Vort/releases)

[MacOSX](https://github.com/donomii/Vort/releases)

## Your own digital safety box

Keep your files securely on a USB key or in the cloud. 

Vort can encrypt your file box, guaranteeing your privacy even if you lose your usb key.  Vort protects your privacy, preventing your cloud host or government agency from reading your files.

## Features

* Undelete
* Compression
* Encryption
* Works well with sync programs like DropBox and CloudMe
* Windows, Linux and MacOSX

# Use

After installing, right click on your desktop, select "New item...", then "Vort".  A new Vort file will appear.  Double click on this to open and use it.

If you don't see a normal file browser window within a few seconds, find the FTP url in the Vort window, and paste that into your web browser.

## Installation

### Windows

Download [Vort](https://github.com/donomii/Vort/releases)

### Linux and MacOSX

Install google's go language, then:

go get github.com/donomii/Vort
go build github.com/donomii/Vort/vort/vort.go


## Command Line Examples

Start Vort with the default options

    vort

Start Vort with encryption

    vort --encrypt=1 --key="a 32-byte key123a 32-byte key123"
    
 connect to Vort with an ftp client
 
     ftp localhost 8021

## Accessing Vort

Vort provides access to your filesystem through an FTP server. Vort will print the url of the server, after it starts up.

You can access this url through normal FTP clients, including:

## Web browsers

Most web browsers include an ftp client. If you have Microsoft Edge, Firefox or Chrome, just paste the vort url into the address bar to access your files.

## Drive mount

Windows has built in support for using FTP servers as ordinary drives. In a file explorer window, click on your computer, then find "Map network drive" somewhere in the menu. Follow the instructions to connect to your file box.

Linux can also mount FTP drives, but it requires installing software and some fiddling about on the command line.

## Stand alone FTP clients

There are a large selection of FTP clients for every platform, and they should all work with Vort.

# Encryption

Warning: Vort uses strong encryption. If you lose your encryption key, you will never be able to access your files again.

Vort uses symmetric key encryption, so in order to access your encrypted filesystem, you must remember the key you created, and type it in every time you open Vort. If you forget or lose your key, all your files are gone forever. They are never coming back. Using encryption is entirely at your own risk, and I am not responsible in any way for any data loss or leak due to using encryption.

# Disclaimer

Actually, I'm not responsible at all for any problems caused by using this software. It is open source, you're getting it for free, and if it doesn't work you're entitled to a full refund and nothing more.
