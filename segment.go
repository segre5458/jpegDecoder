package main

import (
	//"bytes"
	"fmt"
)

func CheckSOI(b int){
	if b != 0xffd8{
		fmt.Println("This file is not JPEG")
	} else{
		fmt.Println("This file is JPEG")
	}
}