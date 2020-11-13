package main

import (
	"fmt"
)

const finLigne = -1
const finImage = -2

func main() {
	var c, d chan int
	var s1 chan bool
	var seuil int

	seuil = 127

	c = make(chan int)
	d = make(chan int)

	go source(c)

	go affichage(d)

	go seuillage(seuil, c, d)
	
	<-s1

	fmt.Println("End")
}

func source(in chan int) {
	image := [12]int{100, 200, 150, finLigne, 32, 250, 18, finLigne, 47, 242, 99, finImage}
	var i int
	i = 0

	for {
		if i == 12 {
			i = 0
		}

		in <- image[i]
		//fmt.Println("Entier envoyé:", image[i])
		i++
	}
}

func affichage(out chan int) {
	var x int
	for {
		x = <-out

		switch {
		case x == finLigne:
			fmt.Println()
		case x == finImage:
			fmt.Println()
			fmt.Println()

		default:
			fmt.Print(x)
		}
	}
}

func seuillage(seuil int, in chan int, out chan int) {
	if seuil > 255 || seuil < 0 {
		fmt.Println("Erreur de valeur du seuil!")
	}

	var valPixel int

	for {
		valPixel = <-in

		switch {
		case valPixel == finLigne:
			out <- valPixel

		case valPixel == finImage:
			out <- valPixel

		default:
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
}
