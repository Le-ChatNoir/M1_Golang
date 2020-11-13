package main

import (
	"fmt"
)

func main() {
	var s1, s2, s3 chan bool
	var c, d chan int
	var seuil int

	seuil = 184

	s1 = make(chan bool)
	s2 = make(chan bool)
	s3 = make(chan bool)

	c = make(chan int)
	d = make(chan int)

	go func() {
		source(c)
		s1 <- true
	}()

	go func() {
		affichage(d)
		s2 <- true
	}()

	go func() {
		seuillage(seuil, c, d)
		s3 <- true
	}()

	<-s1
	<-s2
	<-s3

	fmt.Println("End")
}

func source(in chan int) {
	image := [9]int{100, 200, 150, 32, 250, 18, 47, 242, 99}

	for i := range image {
		in <- image[i]
		//fmt.Println("Entier envoyé:", image[i])
	}
}

func affichage(out chan int) {
	var x int
	for i := 0; i < 9; i++ {
		x = <-out
		fmt.Println("Entier reçu :", x)
	}
}

func seuillage(seuil int, in chan int, out chan int) {
	if seuil > 255 || seuil < 0 {
		fmt.Println("Erreur de valeur du seuil!")
	}

	var valPixel int

	for i := 0; i < 9; i++ {
		valPixel = <-in

		if valPixel < seuil {
			//fmt.Println("Pixel inferieur au seuil. Envoi d'un 0.")
			out <- 0
		}

		if valPixel >= seuil {
			//fmt.Println("Pixel égal ou supérieur au seuil. Envoi d'un 1.")
			out <- 1
		}
	}
}
