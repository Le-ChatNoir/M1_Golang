package main

import (
	"fmt"
)

func main() {
	var c chan bool
	c = make(chan bool)
	
	go func(){
		fmt.Println("Hello")
		c <- true
	}()
	
	<- c
	fmt.Println("world")
}
