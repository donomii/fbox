# fbox

A secure, encrypted file store.

## Download

Warning:  this is experimental software.  Do not use it for anything important.  Make backups often!

[Windows](https://github.com/donomii/fbox/releases)

[Linux](https://github.com/donomii/fbox/releases)

[MacOSX](https://github.com/donomii/fbox/releases)

## Your own digital safety box

Keep your files securely on a USB key. Fbox can encrypt your file box, guaranteeing your privacy even if you lose your usb key.

Keep your files securely in the cloud. Fbox can protect your privacy, preventing your cloud host or government agency from reading your files.

## Features

* Undelete
* Compression
* Encryption
* Works well with sync programs like DropBox and CloudMe
* Windows, Linux and MacOSX

# Use

Download and open fbox.exe. fbox will open a file window showing the contents of your file box, which starts off empty.

fbox will automatically create your file box, called filebox.fbx, in the same directory as fbox.exe. You can move this directory anywhere on your computer, so long as it is in the same directory as fbox.

Double clicking on this file will open fbox to browse the files in this filebox.


## Installation

### Windows

Download [fbox](https://github.com/donomii/fbox/releases)

###Linux and MacOSX

Install google's go language, then:

go get github.com/donomii/fbox
go build github.com/donomii/fbox/fbox.go


## Command Line Examples

Start fbox with the default options

fbox

Start fbox with encryption

fbox --encrypt=1 --key="a 32-byte key123a 32-byte key123"

## Advanced Use

fbox provides access to your filesystem through an FTP server. Fbox will print the url of the server, after it starts up.

You can access this url through normal FTP clients, including:

## Web browsers

Most web browsers include an ftp client. If you have Microsoft Edge, Firefox or Chrome, just paste the url into the address bar to access your files.

## Drive mount

Windows has built in support for using FTP servers as ordinary drives. In a file explorer window, click on your computer, then find "Map network drive" somewhere in the menu. Follow the instructions to connect to your file box.

Linux can also mount FTP drives, but it requires installing software and some fiddling about on the command line.

## Stand alone FTP clients

There are a large selection of FTP clients for every platform, and they should all work with fbox.

# Encryption

Warning: fbox uses strong encryption. If you lose your encryption key, you will never be able to access your files again.

fbox uses symmetric key encryption, so in order to access your encrypted filesystem, you must remember the key you created, and type it in every time you open fbox. If you forget or lose your key, all your files are gone forever. They are never coming back. Using encryption is entirely at your own risk, and I am not responsible in any way for any data loss or leak due to using encryption.

# Disclaimer

Actually, I'm not responsible at all for any problems caused by using this software. It is open source, you're getting it for free, and if it doesn't work you're entitled to a full refund and nothing more.
