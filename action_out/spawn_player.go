package action_out

import (
	"bytes"
	"encoding/gob"
	"log"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
)

type PlayerInitData struct {
	Name              string
	X, Y, Orientation int
	PlaceBlock        []int
}

type NotifyPlayerSpawn struct {
	To_        string
	player     PlayerInitData
	PlaceBlock []int
}

func NewNotifyPlayerSpawn(to string, player PlayerInitData) NotifyPlayerSpawn {
	return NotifyPlayerSpawn{
		To_:    to,
		player: player,
	}
}

func (a NotifyPlayerSpawn) GetData() []byte {

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(a.player)
	if err != nil {
		log.Fatal("connot encode player spawn action :", err)
	}

	prefixedData := append([]byte{action_type.PlayerSpawnId}, buff.Bytes()...)
	return prefixedData

}

func (a NotifyPlayerSpawn) To() string {
	return a.To_
}
