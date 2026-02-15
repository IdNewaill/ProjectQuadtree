package remote_player

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand/v2"
	"net"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_out"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game_action"
)

func handleServerConnection(
	selfName string,
	conn net.Conn,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
	newConnChan chan ConnectionRegister,
	connActionChan chan ConnAction,
) {

	// Demander au jeu d'envoyer worldState
	worldStateChan := make(chan action_type.InitWorldState)
	gameActionChan <- game_action.GetWorldState{Sender: worldStateChan}

	// Quand on le reçois
	worldState := <-worldStateChan

	// Encoder worldState
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(worldState)
	if err != nil {
		fmt.Println("Erreur encoding init world state")
		return
	}

	// Envoyer worldState
	err = sendData(conn, buff.Bytes())
	if err != nil {
		fmt.Println("Erreur sending init world state")
		return
	}

	// Lire les données d'initalisation des clients (playerInitData)
	data, err := readData(conn)
	if err != nil {
		fmt.Println("Erreur reading perr's server worldState")
		return
	}

	// Decoder ces données d'initalisation
	var peerPlayerInitData action_type.PlayerInitData
	buff = *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buff)
	err = dec.Decode(&peerPlayerInitData)
	if err != nil {
		fmt.Println("Erreur decoding perr's server worldState :", err)
		return
	}

	// Prévenir la personne
	if len(worldState.Players) != 0 {

		to := ""
		for i := range worldState.Players {
			to += worldState.Players[i].Name + " "
		}
		to = to[:len(to)-1]

		actionOutChan <- action_out.NewNotifyPlayerSpawn(to, action_out.PlayerInitData{
			Name:        peerPlayerInitData.Name,
			X:           peerPlayerInitData.X,
			Y:           peerPlayerInitData.Y,
			Orientation: peerPlayerInitData.Orientation,
			PlaceBlock:  peerPlayerInitData.PlaceBlock,
		})

	}

	// Faire apparaître le joueur en jeu
	gameActionChan <- game_action.SpawnPlayer{
		Name:        peerPlayerInitData.Name,
		X:           peerPlayerInitData.X,
		Y:           peerPlayerInitData.Y,
		Orientation: peerPlayerInitData.Orientation,
		PlaceBlock:  peerPlayerInitData.PlaceBlock,
	}

	// crée un id pour le joueur
	id := rand.IntN(9223372036854775807)

	// Garder une connexion active avec le joueur pour communiquer sa position, son orientation, etc...
	packerChan := make(chan packet, 10)
	newConnChan <- ConnectionRegister{
		id:          id,
		isServer:    false,
		playerNames: []string{peerPlayerInitData.Name},
		packetChan:  packerChan,
	}

	handleConnection(selfName, conn, id, gameActionChan, actionOutChan, packerChan, connActionChan)
}

func acceptLoop(
	selfName string,
	listener net.Listener,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
	newConnChan chan ConnectionRegister,
	connActionChan chan ConnAction,
) {
	defer listener.Close()

	for {
		// Accepte les connexions
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while acccepting connection :", err)
			continue
		}

		// Handle connection
		go handleServerConnection(selfName, conn, gameActionChan, actionOutChan, newConnChan, connActionChan)
	}
}

func StartServer(
	selfName string,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
	newConnChan chan ConnectionRegister,
	connActionChan chan ConnAction,
) {
	// Cette fonction sert à démarrer un serveur pour communiquer

	// Actuellement, le serveur se trouve sur localhost
	// mais on peut aussi créer un serveur en utilisant son ip actuelle
	// pour permettre à d'autres personnes qui se trouvent sur le même réseau wifi
	// de nous rejoindre

	// J'ai choisi le port 5060 car c'est le port le plus souvent libre d'après ce que j'ai pu voir

	if !configuration.Global.MapFileMode {
		listener, err := net.Listen("tcp", "localhost:5060")
		if err != nil {
			fmt.Println("cannot start server")
			return
		}

		fmt.Println("server started")

		go acceptLoop(selfName, listener, gameActionChan, actionOutChan, newConnChan, connActionChan)
	}
}
