# fbox

Filesystem in a box

A portable filestore.  It is secure, compressed and very convenient.

## Features

* Snapshots (done)
* Compression (done)
* Encryption (done)
* Works well with sync programs (in progress)

## Installation

    go get github.com/donomii/fbox
    go build github.com/donomii/fbox/example

## Example

    example filesystem.box

## Use

fbox provides access to your filesystem through an FTP server.  Fbox will print the url of the server, after it starts up.

You can access this url through normal FTP clients, including:

## Web browsers

Most web browsers include an ftp client.  If you have Microsoft Edge, Firefox or Chrome, just paste the url into the address bar to access your files.

## Drive mount

Windows has built in support for using FTP servers as ordinary drives.  In a file explorer window, click on your computer, then find "mount drive" somewhere in the menu.  Follow the instructions to connect to a network drive.

Linux can also mount FTP drives, but it requires installing software.

## Stand alone FTP clients

There are a large selection of FTP clients for every platform, and they should all work with fbox.

# Encryption

fbox uses symmetric key encryption, so in order to access your encrypted filesystem, you must remember the key you created, and type it in every time you open fbox.  If you forget or lose your key, all your files are gone forever.  They are never coming back.  Using encryption is entirely at your own risk, and I am not responsible in any way for any data loss or leak due to using encryption.

Actually, I'm not responsible at all for any problems caused by using this software.  It is open source, you're getting it for free, and if it doesn't work you're entitled to a full refund.

