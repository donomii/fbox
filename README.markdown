# fbox

Filesystem in a box

## Features

* Snapshots
* Compression
* Encryption
* Works well with sync programs

## Installation

    go get github.com/donomii/fbox
    go build github.com/donomii/fbox/example

## Example

    example filesystem.box

## Use

fbox provides access to your filesystem through an FTP server.  It will print the url to connect to, after it starts up.

You can access this url through normal FTP clients, including:

## Web browsers

Most web browsers include an ftp client.  If you have Microsoft Edge, Firefox or Chrome, just paste the url into the address bar to access your files.

## Drive mount

Windows has built in support for using FTP serveres as ordinary drives.  In a file explorer window, click on your computer, then find "mount drive" somewhere in the menu

Linux can also mount FTP drives, but it requires installing a package to do so.

## Stand alone FTP clients

There are a large selection of FTP clients for every platform, and they should all work with fbox.

