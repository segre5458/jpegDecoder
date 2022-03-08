package main

import(
	"bytes"
)

func Parse(buffer []byte) (err error){
	r := bytes.NewReader(buffer)
	CheckSOI(readBytesAsInt(r,2))
	return err
}