package main

import "fmt"

type Karel string

func (k string)predstavSe(){
	fmt.Println(k)
}

func main() {
	var s Karel
	s = "karel"
	s.predstavSe()
}
