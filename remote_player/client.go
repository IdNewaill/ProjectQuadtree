package remote_player

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand/v2"
	"net"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game_action"
)

func Connect(
	addr string,
	self action_type.PlayerInitData,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
	newConnChan chan ConnectionRegister,
	connActionChan chan ConnAction,
) bool {
	if !configuration.Global.MapFileMode {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("Error connecting to", addr, ":", err)
			return false
		}

		// encoder les données d'initialisation d'un joueur
		var buff bytes.Buffer
		enc := gob.NewEncoder(&buff)
		err = enc.Encode(self)
		if err != nil {
			fmt.Println("Erreur encoding init player init data")
			return false
		}

		// les envoyer
		err = sendData(conn, buff.Bytes())
		if err != nil {
			fmt.Println("Erreur sending init world state")
			return false
		}

		// lire chez le serveur l'état du monde
		data, err := readData(conn)
		if err != nil {
			fmt.Println("Erreur reading server worldState")
			return false
		}

		// decoder la réponse reçue
		var worldState action_type.InitWorldState
		buff = *bytes.NewBuffer(data)
		dec := gob.NewDecoder(&buff)
		err = dec.Decode(&worldState)
		if err != nil {
			fmt.Println("Erreur decoding server worldState :", err)
			return false
		}

		playerNames := make([]string, len(worldState.Players))

		// faire apparaître le joueur en jeu
		for i := range worldState.Players {
			playerNames[i] = worldState.Players[i].Name
			gameActionChan <- game_action.SpawnPlayer{
				Name:        worldState.Players[i].Name,
				X:           worldState.Players[i].X,
				Y:           worldState.Players[i].Y,
				Orientation: worldState.Players[i].Orientation,
			}
		}

		// id attribué au joueur en question
		id := rand.IntN(9223372036854775807)

		// ajouter du pear to pear
		packerChan := make(chan packet, 10)
		newConnChan <- ConnectionRegister{
			id:          id,
			isServer:    true,
			playerNames: playerNames,
			packetChan:  packerChan,
		}

		go handleConnection(self.Name, conn, id, gameActionChan, actionOutChan, packerChan, connActionChan)

		return true
	}

	return true
}
