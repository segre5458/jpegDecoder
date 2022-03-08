package main

import (
	"fmt"
	"os"
	"io/ioutil"
)

func main(){
	fileName := os.Args[1]
	buf, err := ioutil.ReadFile(fileName)
	if err != nil{
		fmt.Println(err)
		return
	}
	err = Parse(buf)
	if err != nil{
		fmt.Println(err)
		return
	}
}