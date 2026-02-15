package game_action

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/character"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
)

type SpawnPlayer struct {
	Name              string
	X, Y, Orientation int
	PlaceBlock        []int
}

func (a SpawnPlayer) Execute(g *game.Game) {
	fmt.Println("player", a.Name, "connected")
	g.OtherPlayers = append(g.OtherPlayers, character.New(a.Name, a.X, a.Y))
}

func decodeSpawnPlayer(data []byte) (SpawnPlayer, error) {

	buff := *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buff)

	action := SpawnPlayer{}
	err := dec.Decode(&action)
	if err != nil {
		fmt.Println("connot decode player_move remote action:", err)
		return action, err
	}

	return action, nil

}
