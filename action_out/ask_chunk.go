package action_out

import (
	"bytes"
	"encoding/gob"
	"log"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
)

type AskChunk struct {
	To_      string
	X, Y     int
	Receiver string
}

func NewAskChunk(to string, x, y int, receiver string) AskChunk {
	return AskChunk{
		To_:      to,
		Receiver: receiver,
		X:        x,
		Y:        y,
	}
}

func (a AskChunk) GetData() []byte {

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(a)
	if err != nil {
		log.Fatal("connot encode ask chunk action :", err)
	}

	prefixedData := append([]byte{action_type.AskChunkActionId}, buff.Bytes()...)
	return prefixedData

}

func (a AskChunk) To() string {
	return a.To_
}
