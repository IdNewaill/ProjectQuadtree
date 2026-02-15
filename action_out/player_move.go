package action_out

import (
	"bytes"
	"encoding/gob"
	"log"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
)

type NotifyPlayerMove struct {
	To_  string
	Name string
	X, Y int
}

func NewNotifyPlayerMove(to string, name string, x, y int) NotifyPlayerMove {
	return NotifyPlayerMove{
		To_:  to,
		Name: name,
		X:    x,
		Y:    y,
	}
}

func (a NotifyPlayerMove) GetData() []byte {

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(a)
	if err != nil {
		log.Fatal("connot encode player move action :", err)
	}

	prefixedData := append([]byte{action_type.PlayerMoveId}, buff.Bytes()...)
	return prefixedData

}

func (a NotifyPlayerMove) To() string {
	return a.To_
}
