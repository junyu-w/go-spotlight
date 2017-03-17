# fileDB

### Introduction

fileDB is a command-line file search engine that does automatic indexing, and supports quries in different format. (Disclaimer: fileDB is actually not a DB, call it a DB because I'm too lazy to make it a service)

And thanks to the power of *Bleve*, fileDB does query analysis before execution, such that the query term you input need not be super accurate (eg. water = watering = watered)

### Install

if you have `go` installed and want to build from source

1. `got get githun.com/DrakeW/fileDB`
2. `cd` into the directory
3. `go build -o fdb`
4. `ln -s $(pwd)/fdb /usr/local/bin/fdb`

otherwise

1. Download the binary release
2. `cd` into the dir of the downloaded binary
3. `chmod 700 ./fdb` to grant user right to execute program
3. `ln -s $(pwd)/fdb /usr/local/bin/fdb`

### How to use

There are two kinds of queries that are supported:

1. *fuzzy query* syntax: `fdb q <whatever word in your mind>`
2. *strict query* syntax: `fdb sq --extension txt --time -3~-1 --words "<make sure you put words inside double quote>"`

More information can be found by running `fdb help` or `fdb help <command>`

### Note

Since `fdb` recursively indexes file information including part of the file content into the index, it takes about 1 second to index 1000 files. Therefore be cautious when you try to index a directory with lots of files. (But the search will be always fast after the indexing phase)