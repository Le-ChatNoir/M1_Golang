package main

import (
	"fmt"
	"sync"
)

const nbChan int = 18
const nbNodes int = 7
const null = -1
const sendInit = 5
const destInit = 1

//structure packet
type packet struct {
	from int
	to int
	num int
}


func main() {

	//Creation du tab de canaux
	var tabChan[nbChan] chan packet

	//Allocation des canaux
	for i := range tabChan {
		tabChan[i] = make(chan packet)
	}

	//Creation de tab bi dimensionnels de canaux. Supports des tabs de canaux I/O
	var comIn[nbNodes] [] chan packet
	var comOut[nbNodes] [] chan packet

	//Definition des chans I/O pour chque noeud
	//L'ordre des ports doit etre respecte
	comIn[0] = [] chan packet {tabChan[1], tabChan[4], tabChan[9], tabChan[7]}
	comOut[0] = [] chan packet {tabChan[0], tabChan[5], tabChan[8], tabChan[6]}
	comIn[1] = [] chan packet {tabChan[0], tabChan[3]}
	comOut[1] = [] chan packet {tabChan[1], tabChan[2]}
	comIn[2] = [] chan packet {tabChan[2], tabChan[5]}
	comOut[2] = [] chan packet {tabChan[3], tabChan[4]}
	comIn[3] = [] chan packet {tabChan[8], tabChan[11]}
	comOut[3] = [] chan packet {tabChan[9], tabChan[10]}
	comIn[4] = [] chan packet {tabChan[6], tabChan[10], tabChan[13], tabChan[16]}
	comOut[4] = [] chan packet {tabChan[7], tabChan[11], tabChan[12], tabChan[17]}
	comIn[5] = [] chan packet {tabChan[12], tabChan[15]}
	comOut[5] = [] chan packet {tabChan[13], tabChan[14]}
	comIn[6] = [] chan packet {tabChan[17], tabChan[14]}
	comOut[6] = [] chan packet {tabChan[16], tabChan[15]}

	//Declaration du waitGroup
	var wg sync.WaitGroup

	//Nombre de synchros a attendre
	wg.Add(nbNodes)

	//Lancement des goroutines
	for i := 0; i < nbNodes; i++ {
		go func(i int){
			node(i, comIn[i], comOut[i], i == sendInit, destInit)
			//Decremente le nombre de sync en attente
			wg.Done()
		} (i)
	}

	//Attente de la completion du waitGroup. Ok quand valeur du add == nombre de Done
	wg.Wait()

	fmt.Println("End")
}

func node(id int, in, out []chan packet, innondation bool, idDestination int){
	//Declaration de messIn et messOut selon le nombre de canaux
	var messOut []packet
	var messIn []packet
	var history []int
	history = make([]int, nbNodes)
	messIn = make([]packet, len(in))
	messOut = make([]packet, len(out))

	var p packet

	//Remplissage de messOut avec le packet si noeud qui propage choisi
	if(innondation){
		p = packet{id, idDestination, 0}
	//Si noeud non propagateur, envoie des packets null
	} else {
		p = packet{null, null, null}
	}

	for i := range out{
		messOut[i] = p
	}

	for i := range history{
		history[i] = null
	}

	//Debut boucle comm infinie
	for tour := 0; true; tour++ {
		communication(in, messIn, out, messOut)

		p = packet{null, null, null}
		//Examination des messages recus
		for i := range in{
			//Si pas un message null, aka si un packet arrive
			if messIn[i].from != null {
				//Si packet pour moi, je l'affiche et ne renvoie que des nulls
				if messIn[i].to == id{
					fmt.Println("Noeud", id, "recoit message de", messIn[i].from, "tour", messIn[i].num)
				//Sinon, je passe mes propres messages pour transmettre le packet a mes voisins
				} else if messIn[i].num > history[messIn[i].from] {
					p = messIn[i]
					history[messIn[i].from] = messIn[i].num
				}
			}
		}
		for i := range out{
			if messIn[i].from != null{
				messOut[i] = packet{null, null, null}
			} else {
				messOut[i] = p
			}
		}

	}
	//Fin boucle infinie

}

func receive(messIn []packet, in []chan packet, wg *sync.WaitGroup){
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

func send(messOut []packet, out []chan packet, wg *sync.WaitGroup){
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

func communication(in []chan packet, messIn []packet, out []chan packet, messOut []packet){
	var wg sync.WaitGroup

	receive(messIn, in, &wg )
	send(messOut, out, &wg )

	wg.Wait()
}