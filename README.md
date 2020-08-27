## Overview of tools

* [gounzip](cmd/gounzip) - is a tool to unzip aes encrypted archives
* [gozip](cmd/gozip) - is a tool to create aes encrypted archives
* [mwcompress](cmd/mwcompress) - MalWare Compress - is a tool to pack malware into an aes encrypted zip archive
* [mwextract](cmd/mwextract) - MalWare Extract - is a tool to unpack malware from aes enrypted zip archives

### gozip

#### Syntax
```
gozip <ARCHIVE.ZIP> <PASSWORD> <File(s)>...
```

### gounzip

#### Syntax

```
gounzip -f <ARCHIVE.ZIP> -p <PASSWORD>
```

### mwcompress

`mwcompress` by default uses the industry standard password `infected`.

#### Syntax

In the most simple case, useing default settings (including default password) you call `mwcompress` with 2 arguments.

```
mwcompress <ARCHIVE.ZIP> <File(s)>...
```

Off course you can set a custom password, too.

```
mwcompress -p <PASSWORD> <ARCHIVE.ZIP> <File(s)>...
```

### mwextract

`mwextract` is the opposite of `mwcompress`, it unpacks files from a aes encrypted zip archive.
The industry standard password (`infected`) is used if no other password is set.

#### Syntax

```
mwextract <ARCHIVE.ZIP>...
```

### build 

```
fwa@fwa01lt:~/go/src/github.com/scusi/AesZip$ docker run -ti --rm -v `pwd`:/usr/src/aeszip -w /usr/src/aeszip golang:1.15 bash
root@7828a19dc519:/usr/src/aeszip# go build ./cmd/mwcompress/
go: downloading github.com/alexmullins/zip v0.0.0-20180717182244-4affb64b04d0
go: downloading golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
root@7828a19dc519:/usr/src/aeszip# go build ./cmd/mwextract/ 
go: downloading github.com/scusi/MultiChecksum v0.1.0
go: downloading github.com/dchest/blake2b v1.0.0
go: downloading github.com/dchest/blake2s v1.0.0
root@7828a19dc519:/usr/src/aeszip# go build ./cmd/gounzip/  
root@7828a19dc519:/usr/src/aeszip# go build ./cmd/gozip/  
root@7828a19dc519:/usr/src/aeszip# exit
fwa@fwa01lt:~/go/src/github.com/scusi/AesZip$ 
```

