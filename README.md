## Overview of tools

* [gounzip](cmd/gounzip) - is a tool to unzip aes encrypted archives
* [gozip](cmd/gozip) - is a tool to create aes encrypted archives
* [mwcompress](cmd/mwcompress) - MalWare Compress - is a tool to pack malware into an aes encrypted zip archive
* [mwextract](cmd/mwextract) - MalWare Extract - is a tool to unpack malware from aes enrypted zip archives

### gozip

#### build

```
go build ./cmd/gozip
```

#### Syntax
```
gozip <ARCHIVE.ZIP> <PASSWORD> <File(s)>...
```

### gounzip

#### build

```
go build ./cmd/gounzip
```

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


