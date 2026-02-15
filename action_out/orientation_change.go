package action_out

import (
	"bytes"
	"encoding/gob"
	"log"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
)

type NotifyOrientationChange struct {
	To_            string
	PlayerName     string
	NewOrientation int
}

func NewNotifyOrientationChange(to string, playerName string, newOrientation int) NotifyOrientationChange {
	return NotifyOrientationChange{
		To_:            to,
		PlayerName:     playerName,
		NewOrientation: newOrientation,
	}
}

func (a NotifyOrientationChange) GetData() []byte {

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(a)
	if err != nil {
		log.Fatal("connot encode orientation change action :", err)
	}

	prefixedData := append([]byte{action_type.OrientationChangeId}, buff.Bytes()...)
	return prefixedData

}

func (a NotifyOrientationChange) To() string {
	return a.To_
}
