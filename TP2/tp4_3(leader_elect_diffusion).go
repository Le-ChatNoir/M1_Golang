package main

import (
	"fmt"
)

//Il faut nbNodes == nbTours
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

	//Tableau de l'ordre d'importance des noeuds modifiable librement
	nodeNumber := [4]int{78, 54, 14, 32}
	//Lancement d'autant de goroutine que de noeuds. Utilisation de k%nbNodes pour que le chanel de sortie au 4eme tour (qui sera alors k+1 = 5) redevienne de channel d'entree 1 pour faire une boucle
	for k := 0; k < nbNodes; k++ {
		go func(k int){
			node (k,tabCan[k],tabCan[(k+1)%nbNodes], nodeNumber[k])

			sync[k] <-true
		}(k)
	}
	
	//Vidage des canneaux de synchro quand la tache est effectuee
	for i:= range sync {
		<- sync[i]
	}

	fmt.Println("End")
}



func node(id int, in, out chan int, nodeNumber int){
	var messIn, messOut int
	var synchro chan bool
	var isLeader bool = false
	var leaderID int = -1

	messOut = nodeNumber

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

		//Fin communication

		fmt.Println("Noeud numero ", id, ", (mon id unique : ", nodeNumber, ")", "Leader: ", isLeader)
		//Si le processus a le jeton, alors il setup son prochain message a transmettre pour etre le jeton, sinon, le prochain message est null
		switch {
		case messIn == null:
			fmt.Println("	En tete (recu: ", messIn, ")")
			messOut = null
		case messIn < nodeNumber:
			fmt.Println("	En tete (recu: ", messIn, ")")
			messOut = null
		case messIn > nodeNumber:
			fmt.Println("	Depassement (recu: ", messIn, ")")
			messOut = messIn
		case messIn == nodeNumber:
			isLeader = true
			fmt.Println("	Je suis le leader! (recu: ", messIn, ")")
		}

	}

	fmt.Println("\nEnvoi de l'ID leader\n")
	//Communication
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

		//Fin communication

		//Si le processus est le leader, il envoie son ID, sinon, il se met en attente de l'id du leader
		if leaderID == -1 {
			if isLeader {
				fmt.Println("	Je suis le leader. Envoi de mon ID")
				messOut = id
				leaderID = id
			} else {
				//Si il n'est pas le leader, et qu'il ne recoit pas les messages poubelle null des precedentes boucles, alors il met a jour son leader
				if messIn != null {
					fmt.Println("	Noeud ", id, " a recu l'ID du leader (", messIn, ")")
					messOut = messIn
					leaderID = messIn
				}
			}
		}
	}
	fmt.Println("Je suis le noeud numero ", id, " et le leader est le noeud numero ", leaderID)
}