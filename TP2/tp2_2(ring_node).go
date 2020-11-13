package main

import (
	"fmt"
)

const nbNodes int = 4
const nbTours int = 4
const null int = -1

func main() {
	var sync[nbNodes] chan bool
	var tabCan [nbNodes] chan int

	//Boucle de make pour gain de temps
	for i := range tabCan {
		tabCan[i] = make(chan int)
	}
	//Initialisation de tous les cannaux de synchronisation (un pour chaque goroutine lancee)
	//Alternativement, on peut utiliser un seul canal avec sync <- true l.31 et une boucle for i < nbNodes <- sync l.36. Le canal agira alors come une pile
	for i:= range sync {
		sync[i] = make(chan bool)
	}


	//Lancement d'autant de goroutine que de noeuds. Utilisation de k%nbNodes pour que le chanel de sortie au 4eme tour (qui sera alors k+1 = 5) redevienne de channel d'entree 1 pour faire une boucle
	for k := 0; k < nbNodes; k++ {
		go func(k int){
			node (k,tabCan[k],tabCan[(k+1)%nbNodes])

			sync[k] <-true
		}(k)
	}
	
	//Vidage des canneaux de synchro quand la tache est effectuee
	for i:= range sync {
		<- sync[i]
	}

	fmt.Println("End")
}



func node(id int, in, out chan int){
	var x int
	var synchroInt1, synchroInt2 chan bool

	//Canneaux de synchro interne car sinon x = <- fait un deadlock, il faut donc que le recu et l'envoi se lancent en parallele
	synchroInt1 = make(chan bool)
	synchroInt2 = make(chan bool)

	//Plusieurs tous sont effecutes. on peut se contenter de synchroInt1 et synchroInt2 car ils sont vides puis re-replis a chaque tour
	for j:=0; j<nbTours;j++{
		
		go func(){
			x = <- in
			synchroInt1 <- true
		}()

		go func(){
			out <- null
			synchroInt2 <- true
		}()

		<- synchroInt1
		<- synchroInt2

		fmt.Println("Je suis le noeud numero", id, ", j'ai recu l'entier", x, "et j'en suis au tour", j)
	}
}