package main

import (
	"fmt"
	"sync"
)

const nbChan int = 18
const nbNodes int = 7

func main() {
	//Creation du tab de canaux
	var tabChan[nbChan] chan int

	//Allocation des canaux
	for i := range tabChan {
		tabChan[i] = make(chan int)
	}

	//Creation de tab bi dimensionnels de canaux. Supports des tabs de canaux I/O
	var comIn[nbNodes] [] chan int
	var comOut[nbNodes] [] chan int

	//Definition des chans I/O pour chque noeud
	//L'ordre des ports doit etre respecte
	comIn[0] = [] chan int {tabChan[1], tabChan[4], tabChan[9], tabChan[7]}
	comOut[0] = [] chan int {tabChan[0], tabChan[5], tabChan[8], tabChan[6]}
	comIn[1] = [] chan int {tabChan[0], tabChan[3]}
	comOut[1] = [] chan int {tabChan[1], tabChan[2]}
	comIn[2] = [] chan int {tabChan[2], tabChan[5]}
	comOut[2] = [] chan int {tabChan[3], tabChan[4]}
	comIn[3] = [] chan int {tabChan[8], tabChan[11]}
	comOut[3] = [] chan int {tabChan[9], tabChan[10]}
	comIn[4] = [] chan int {tabChan[6], tabChan[10], tabChan[13], tabChan[16]}
	comOut[4] = [] chan int {tabChan[7], tabChan[11], tabChan[12], tabChan[17]}
	comIn[5] = [] chan int {tabChan[12], tabChan[15]}
	comOut[5] = [] chan int {tabChan[13], tabChan[14]}
	comIn[6] = [] chan int {tabChan[17], tabChan[14]}
	comOut[6] = [] chan int {tabChan[16], tabChan[15]}

	//Declaration du waitGroup
	var wg sync.WaitGroup

	//Nombre de synchros a attendre
	wg.Add(nbNodes)

	//Lancement des goroutines
	for i := 0; i < nbNodes; i++ {
		go func(i int){
			node(i, comIn[i], comOut[i])
			//Decremente le nombre de sync en attente
			wg.Done()
		} (i)
	}

	//Attente de la completion du waitGroup. Ok quand valeur du add == nombre de Done
	wg.Wait()

	fmt.Println("End")
}

func node(id int, in, out []chan int){
	//Declaration de messIn et messOut selon le nombre de canaux
	var messOut []int
	var messIn []int
	messIn = make([]int, len(in))
	messOut = make([]int, len(out))

	//Remplissage de messOut aved son id de noeud pour l'envoyer a chaque voisin
	for i := range out{
		messOut[i] = id
	}

	communication(in, messIn, out, messOut)

	//Possibilite de print un tableau deja forme
	fmt.Println("Voisins du noeud ", id, ":", messIn)
}

func receive(messIn []int, in []chan int, wg *sync.WaitGroup){
	var nbChans int
	nbChans = len(in)

	wg.Add(nbChans)

	for i := 0; i < nbChans; i++ {
		go func(i int){
			//Reception du message dans messIn
			messIn[i] = <- in[i]
			//Decremente le nombre de sync en attente
			wg.Done()
		} (i)
	}
}

func send(messOut []int, out []chan int, wg *sync.WaitGroup){
	var nbChans int
	nbChans = len(out)

	wg.Add(nbChans)

	for i := 0; i < nbChans; i++ {
		go func(i int){
			//Envoi du message de messOut dans out
			out[i] <- messOut[i]
			//Decremente le nombre de sync en attente
			wg.Done()
		} (i)
	}
}

func communication(in []chan int, messIn []int, out []chan int, messOut []int){
	var wg sync.WaitGroup

	send(messOut, out, &wg )
	receive(messIn, in, &wg )

	wg.Wait()
}