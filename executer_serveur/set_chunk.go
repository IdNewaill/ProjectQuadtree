package action_out

import (
	"bytes"
	"encoding/gob"
	"log"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
)

type SetChunk struct {
	To_     string
	X, Y    int
	Content [][]int
}

func NewSetChunk(to string, x, y int, content [][]int) SetChunk {
	return SetChunk{
		To_:     to,
		X:       x,
		Y:       y,
		Content: content,
	}
}

func (a SetChunk) GetData() []byte {

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(a)
	if err != nil {
		log.Fatal("connot encode ask chunk action :", err)
	}

	prefixedData := append([]byte{action_type.SetChunkActionId}, buff.Bytes()...)
	return prefixedData

}

func (a SetChunk) To() string {
	return a.To_
}
