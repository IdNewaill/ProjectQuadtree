package remote_player

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game_action"
)

type relayActionOut struct {
	Data_ []byte
	To_   string
}

func (a relayActionOut) GetData() []byte { return a.Data_ }
func (a relayActionOut) To() string      { return a.To_ }

func sendData(conn net.Conn, data []byte) error {
	// Envoyer les données reçues

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(len(data)))

	_, err := conn.Write(buf)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

func readData(conn net.Conn) ([]byte, error) {
	// Cette fonction sert à les données reçues

	lengthBuf := make([]byte, 8)

	_, err := io.ReadFull(conn, lengthBuf)
	if err != nil {
		return nil, err
	}

	packetLen := binary.BigEndian.Uint64(lengthBuf)

	dataBuf := make([]byte, packetLen)
	_, err = io.ReadFull(conn, dataBuf)
	if err != nil {
		return nil, err
	}

	return dataBuf, nil
}

func handleConnection(
	selfName string,
	conn net.Conn,
	peerId int,
	gameActionChan chan game.GameAction,
	actionOutChan chan game.ActionOut,
	packetChan chan packet,
	connActionChan chan ConnAction,
) {
	defer conn.Close()

	// Cette fonction sert à envoyer des données à un joueur et communiquer avec lui
	go func() {
		for {
			currentPacket := <-packetChan

			//Incorporer "to" dans le paquet
			toInByte := []byte(currentPacket.to)
			buf := make([]byte, 8)
			binary.BigEndian.PutUint64(buf, uint64(len(toInByte)))
			buf = append(buf, toInByte...)
			buf = append(buf, currentPacket.data...)

			// envoyer paquet
			err := sendData(conn, buf)
			if err != nil {
				fmt.Println("Error while sending data")

				// arrête l'envoie à cette personne
				// tout le monde sera prévenu que ce joueur ne se trouve plus dans le jeu
				connActionChan <- ConnAction{
					id:     peerId,
					action: false,
					name:   "all",
				}

				return
			}
		}
	}()

	// lire les données reçues
	for {
		data, err := readData(conn)
		if err != nil || len(data) < 8 {

			// arrête l'envoie à cette personne
			// tout le monde sera prévenu que ce joueur ne se trouve plus dans le jeu
			connActionChan <- ConnAction{
				id:     peerId,
				action: false,
				name:   "all",
			}

			return
		}

		// decoder la destination du paquet reçu
		toLen := binary.BigEndian.Uint64(data[:8])
		to := string(data[8 : 8+toLen])
		packet := data[8+toLen:]

		// transmettre paquet
		if to != configuration.Global.PlayerName && to != "up" {
			fmt.Println("relay packet to", to)
			actionOutChan <- relayActionOut{
				Data_: packet,
				To_:   to,
			}
		}

		// ne pas executer une action si elle ne me concerne pas
		fmt.Println(to)
		if to == "up" || to == "all" || strings.Contains(to, selfName) {
			remoteAction, err := game_action.GetRemoteActionFrom(packet)
			if err != nil {
				fmt.Println("Error while parsing peer packet :", err)
				continue
			}

			if playerSpawnAction, ok := remoteAction.(game_action.SpawnPlayer); ok {
				connActionChan <- ConnAction{
					id:     peerId,
					action: true,
					name:   playerSpawnAction.Name,
				}
			}

			if playerDespawnAction, ok := remoteAction.(game_action.DespawnPlayer); ok {
				connActionChan <- ConnAction{
					id:     peerId,
					action: false,
					name:   playerDespawnAction.Name,
				}
			}

			gameActionChan <- remoteAction
		}
	}
}
