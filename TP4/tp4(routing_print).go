package main

import (
	"fmt"
	"sync"
)

const nbChan int = 18
const nbNodes int = 7
const null = -1
const infini int = 16

//structure [nbNodes] infoDest
type infoDest struct {
	//Indice du canal de sortie vers lequel il faut diriger le message au noeud x
	indVois int
	//nb sauts pour atteindre le noeud x
	cout int
}


func main() {

	//Creation du tab de canaux
	var tabChan[nbChan] chan [nbNodes] infoDest

	//Allocation des canaux
	for i := range tabChan {
		tabChan[i] = make(chan [nbNodes] infoDest)
	}

	//Creation de tab bi dimensionnels de canaux. Supports des tabs de canaux I/O
	var comIn[nbNodes] [] chan [nbNodes] infoDest
	var comOut[nbNodes] [] chan [nbNodes] infoDest

	//Definition des chans I/O pour chque noeud
	//L'ordre des ports doit etre respecte
	comIn[0] = [] chan [nbNodes] infoDest {tabChan[1], tabChan[4], tabChan[9], tabChan[7]}
	comOut[0] = [] chan [nbNodes] infoDest {tabChan[0], tabChan[5], tabChan[8], tabChan[6]}
	comIn[1] = [] chan [nbNodes] infoDest {tabChan[0], tabChan[3]}
	comOut[1] = [] chan [nbNodes] infoDest {tabChan[1], tabChan[2]}
	comIn[2] = [] chan [nbNodes] infoDest {tabChan[2], tabChan[5]}
	comOut[2] = [] chan [nbNodes] infoDest {tabChan[3], tabChan[4]}
	comIn[3] = [] chan [nbNodes] infoDest {tabChan[8], tabChan[11]}
	comOut[3] = [] chan [nbNodes] infoDest {tabChan[9], tabChan[10]}
	comIn[4] = [] chan [nbNodes] infoDest {tabChan[6], tabChan[10], tabChan[13], tabChan[16]}
	comOut[4] = [] chan [nbNodes] infoDest {tabChan[7], tabChan[11], tabChan[12], tabChan[17]}
	comIn[5] = [] chan [nbNodes] infoDest {tabChan[12], tabChan[15]}
	comOut[5] = [] chan [nbNodes] infoDest {tabChan[13], tabChan[14]}
	comIn[6] = [] chan [nbNodes] infoDest {tabChan[17], tabChan[14]}
	comOut[6] = [] chan [nbNodes] infoDest {tabChan[16], tabChan[15]}

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

func node(id int, in, out []chan [nbNodes] infoDest){
	//Declaration de tablesIn et tablesOut selon le nombre de canaux
	var tablesOut [][nbNodes] infoDest
	var tablesIn [][nbNodes] infoDest
	var tableRoutage [nbNodes] infoDest
	var modif bool
	var nbModif int = 0
	tablesIn = make([][nbNodes] infoDest, len(in))
	tablesOut = make([][nbNodes] infoDest, len(out))

	for i := range tableRoutage{
		tableRoutage[i] = infoDest{-1, infini}
	}

	tableRoutage[id].cout = 0

	for i := range out {
		tablesOut[i] = tableRoutage
	}

	//Debut boucle comm infinie
	for tour := 0; true; tour++ {
		modif = false
		communication(in, tablesIn, out, tablesOut)
		//tableIn continet les tables de routage de tous les voisins
		//Donc il faut faire une double boucle pour trouver le port (case de tableIn)
		//Et une autre pour mettre à jour les tableRoutage

		//I correspond au port par lequel l'info arrive
		for i := range in{
			//Et j case du tableau tableRoutage
			for j := range tableRoutage{
				//Si le nouveau chemin est plus court que le precedent, on remplace
				if tablesIn[i][j].cout + 1 < tableRoutage[j].cout {
					tableRoutage[j].cout = tablesIn[i][j].cout +1
					tableRoutage[j].indVois = i
					modif = true
				}
			}
		}
		
		//Affichage
		if modif{
			nbModif++
			fmt.Println("id:", id, "table routage:")
			for i := range tableRoutage{
				fmt.Println("		Noeud", i, tableRoutage[i])
			}
			fmt.Println("nbModif:", nbModif, "\n")
		}

		//Remise a zero par securite
		for i := range out {
			tablesOut[i] = tableRoutage
		}
	}
	//Fin boucle infinie
}

func receive(tablesIn [][nbNodes] infoDest, in []chan [nbNodes] infoDest, wg *sync.WaitGroup){
	var nbChans int
	nbChans = len(in)

	wg.Add(nbChans)

	for i := 0; i < nbChans; i++ {
		go func(i int){
			//Reception du message dans tablesIn
			tablesIn[i] = <- in[i]
			//Decremente le nombre de sync en attente
			wg.Done()
		} (i)
	}
}

func send(tablesOut [][nbNodes] infoDest, out []chan [nbNodes] infoDest, wg *sync.WaitGroup){
	var nbChans int
	nbChans = len(out)

	wg.Add(nbChans)

	for i := 0; i < nbChans; i++ {
		go func(i int){
			//Envoi du message de tablesOut dans out
			out[i] <- tablesOut[i]
			//Decremente le nombre de sync en attente
			wg.Done()
		} (i)
	}
}

func communication(in []chan [nbNodes] infoDest, tablesIn [][nbNodes] infoDest, out []chan [nbNodes] infoDest, tablesOut [][nbNodes] infoDest){
	var wg sync.WaitGroup

	receive(tablesIn, in, &wg )
	send(tablesOut, out, &wg )

	wg.Wait()
}