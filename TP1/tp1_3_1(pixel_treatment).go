package main

import (
	"fmt"
)

func main() {
	var s1, s2, s3 chan bool
	var c, d chan int
	var seuil, pixel int

	seuil = 184
	pixel = 245

	s1 = make(chan bool)
	s2 = make(chan bool)
	s3 = make(chan bool)

	c = make(chan int)
	d = make(chan int)

	go func() {
		source(c, pixel)
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

func source(in chan int, x int) {
	in <- x
	fmt.Println("Entier envoyé:", x)
}

func affichage(out chan int) {
	var x int
	x = <-out
	fmt.Println("Entier reçu :", x)
}

func seuillage(seuil int, in chan int, out chan int) {
	if seuil > 255 || seuil < 0 {
		fmt.Println("Erreur de valeur du seuil!")
	}

	var valPixel int
	valPixel = <-in

	if valPixel < seuil {
		fmt.Println("Pixel inferieur au seuil. Envoi d'un 0.")
		out <- 0
	}

	if valPixel >= seuil {
		fmt.Println("Pixel égal ou supérieur au seuil. Envoi d'un 1.")
		out <- 1
	}
}
