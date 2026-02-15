package game_action

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
)

type PlayerMove struct {
	Name string
	X, Y int
}

func (a PlayerMove) Execute(g *game.Game) {

	for i := range g.OtherPlayers {
		if g.OtherPlayers[i].Name == a.Name {
			g.OtherPlayers[i].RemoteMove(a.X, a.Y)
			return
		}
	}

	fmt.Println("Error : player", a.Name)
}

func decodePlayerMove(data []byte) (PlayerMove, error) {
	
	action := PlayerMove{}

	buff := *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&action)
	if err != nil {
		fmt.Println("connot decode player_move remote action:", err)
		return action, err
	}

	return action, nil

}
