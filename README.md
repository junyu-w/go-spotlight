# Go Spotlight

### Introduction

Go Spotlight is a command-line file search engine that does automatic indexing, and supports quries in different format.

And thanks to the power of *bleve*, `gsp` does query analysis before execution, such that the query term you input need not be super accurate (eg. water = watering = watered)

### Install

if you have `go` installed and want to build from source

1. `got get githun.com/DrakeW/go-spotlight`
2. `cd` into the directory
3. `go build -o gsp`
4. `ln -s $(pwd)/fdb /usr/local/bin/gsp`

otherwise

1. Download the binary release
2. `cd` into the dir of the downloaded binary
3. `chmod 700 ./gsp` to grant user right to execute program
3. `ln -s $(pwd)/fdb /usr/local/bin/gsp`

### How to use

There are two kinds of queries that are supported:

1. *fuzzy query* syntax: `gsp q <whatever word in your mind>`
2. *strict query* syntax: `gsp sq --extension txt --time -3~-1 --words "<make sure you put words inside double quote>"`

More information can be found by running `gsp help` or `gsp help <command>`

### Note

Since `gsp` recursively indexes file information including part of the file content into the index, it takes about 1 second to index 1000 files. Therefore be cautious when you try to index a directory with lots of files. (But the search will be always fast after the indexing phase)