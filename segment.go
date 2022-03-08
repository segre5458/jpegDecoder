package main

import (
	//"bytes"
	"fmt"
)

func CheckSOI(b int){
	if b != 0xffd8{
		fmt.Println("This file is not JPEG")
	}
}

func CheckEOI(b int)(f bool){
	if b != 0xffd9{
		f = false
	} else{
		f = true
	}

	return f
}