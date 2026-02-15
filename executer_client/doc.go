package remote_player

import (
	"strings"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_out"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game_action"
)

type packet struct {
	to   string
	data []byte
}

type ConnAction struct {
	id     int
	action bool
	name   string // il est biensur vérifier que le nom du joueur n'est pas égal à "all" car "all" permet d'envoyer à tout le monde un paquet
}

type ConnectionRegister struct {
	id          int
	isServer    bool
	packetChan  chan packet
	playerNames []string
}

func inCommon(a1, a2 []string) (common []string) {
	for i1 := range a1 {
		for i2 := range a2 {
			if a1[i1] == a2[i2] {
				common = append(common, a1[i1])
			}
		}
	}
	return
}

func managePeers(
	newConnChan chan ConnectionRegister,
	connActionChan chan ConnAction,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
) {
	peers := []ConnectionRegister{}

	for {
		select {
		case conn := <-newConnChan:
			peers = append(peers, conn)

		case connAction := <-connActionChan:
			for peerIndex := range peers {
				if peers[peerIndex].id == connAction.id {

					if connAction.action {

						peers[peerIndex].playerNames = append(peers[peerIndex].playerNames, connAction.name)

					} else {

						// connexion fermée
						if connAction.name == "all" {

							// retirer le joueur de tous les peers
							for i := range peers[peerIndex].playerNames {

								// prévenir tout le monde
								actionOutChan <- action_out.NewNotifyPlayerDespawn(
									"all",
									peers[peerIndex].playerNames[i],
								)

								// retirer le joueur du jeu
								gameActionChan <- game_action.DespawnPlayer{
									Name: peers[peerIndex].playerNames[i],
								}
							}

							// remove peer
							peers[peerIndex] = peers[len(peers)-1]
							peers = peers[:len(peers)-1]

							break
						}

						// retirer le joueur
						for i := range peers[peerIndex].playerNames {

							if peers[peerIndex].playerNames[i] == connAction.name {

								namesLen := len(peers[peerIndex].playerNames)
								peers[peerIndex].playerNames[i] = peers[peerIndex].playerNames[namesLen-1]
								peers[peerIndex].playerNames = peers[peerIndex].playerNames[:namesLen-1]

								break
							}

						}
					}
				}
			}

		case actionOut := <-actionOutChan:

			data := actionOut.GetData()
			to := actionOut.To()

			for peerIndex := range peers {

				//si peerTo est vide, c'est qu'il n'y a rien à envoyer
				peerTo := ""
				if to == "up" {

					if peers[peerIndex].isServer {
						peerTo = "up"
					}

				} else {

					peerToNames := peers[peerIndex].playerNames
					if to != "all" {
						peerToNames = inCommon(
							strings.Split(to, " "),
							peers[peerIndex].playerNames,
						)
					}

					if len(peerToNames) == 0 {
						continue
					}

					for i := range peerToNames {
						peerTo += peerToNames[i] + " "
					}

					peerTo = peerTo[:len(peerTo)-1]

				}
				if len(peerTo) == 0 {
					continue
				}

				peers[peerIndex].packetChan <- packet{
					to:   peerTo,
					data: data,
				}
			}
		}
	}
}

func OpenConnections(
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
) (chan ConnectionRegister, chan ConnAction) {
	// Cette fonction s'occupe de gérer les nouvelles connexion de clients

	newConnChan := make(chan ConnectionRegister, 10)
	connActionChan := make(chan ConnAction, 10)

	go managePeers(newConnChan, connActionChan, gameActionChan, actionOutChan)

	return newConnChan, connActionChan
}
