# Web editor 

Written in Go using websockets and Ace editor

- Ace editor - https://ace.c9.io/

## Features

- edit html, css & js files
- multiple users edit same file near real time
- file contents are loaded to app memory until all users are ready to save new contents
- when every user is disconnect, file is saved and unloaded from memory
- when multiple users are active, every user must set Ready state
- file can be loaded from file system or http
- app can keep versions in separate folder

## Environment variables

required:

- FILE_IO - file or http 
- FILE_PATH - path to file

```
FILE_IO: "file"
FILE_PATH: "./assets/custom.css"
```

optional:

-  app can save each version in separate folder/location
- VERSIONS_IO - file or http 
- VERSIONS_PATH - path to directory

```
VERSIONS_IO: "file"
VERSIONS_PATH: "./assets/versions/"
```

- user connection ttl 
- CONN_TTL - `1m/1h/1d`

```
CONN_TTL: "1h"
```

## Run

- build docker image

```
make docker-build
```

- run app with docker-compose
```
make compose-start
```

- stop app with docker-compose
```
make compose-remove
```

## Ideas for future

- notification for other users to save
 - reminder every x mins

- read dir & open multiple files

- api
- jwt login