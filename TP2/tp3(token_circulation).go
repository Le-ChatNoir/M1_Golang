package main

import (
	"fmt"
)

const nbNodes int = 4
const nbTours int = 4
const null int = -1
const jeton int = -2

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
			node (k,tabCan[k],tabCan[(k+1)%nbNodes], k==1)

			sync[k] <-true
		}(k)
	}
	
	//Vidage des canneaux de synchro quand la tache est effectuee
	for i:= range sync {
		<- sync[i]
	}

	fmt.Println("End")
}



func node(id int, in, out chan int, gotToken bool){
	var messIn, messOut int
	var synchro chan bool

	//On setup le 1er message a envoyer si le processus node lance est celui detenteur du jeton
	if gotToken {
		messOut = jeton
	} else {
		messOut = null
	}

	//Canal de synchro interne car sinon x = <- fait un deadlock, il faut donc que le recu et l'envoi se lancent en parallele
	synchro = make(chan bool)

	//PARTIE COMMUNICATION: Ne pas toucher
	//Plusieurs tous sont effecutes. on peut se contenter de synchroInt1 et synchroInt2 car ils sont vides puis re-replis a chaque tour
	for j:=0; j<nbTours;j++{

		go func(){
			messIn = <- in
			synchro <- true
		}()

		go func(){
			out <- messOut
			synchro <- true
		}()

		<- synchro
		<- synchro

		fmt.Println("Je suis le noeud numero ", id)
		//Si le processus a le jeton, alors il setup son prochain message a transmettre pour etre le jeton, sinon, le prochain message est null
		if messIn == jeton {
			fmt.Println("		Et j'ai le token!")
			messOut = jeton
		} else {
			messOut = 0
		}
	}
}