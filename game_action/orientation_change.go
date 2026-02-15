package game_action

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
)

type RemoteOrientationChange struct {
	PlayerName string
	NewOrientation int
}

func (a RemoteOrientationChange) Execute(g *game.Game) {

	for i := range g.OtherPlayers {
		if g.OtherPlayers[i].Name == a.PlayerName {
			g.OtherPlayers[i].RemoteSetOrientation(a.NewOrientation)
			return
		}
	}

	fmt.Println("Error : player", a.PlayerName, "not found")
}

func decodeOrientationChange(data []byte) (RemoteOrientationChange, error) {

	action := RemoteOrientationChange{}

	buff := *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&action)
	if err != nil {
		fmt.Println("connot decode orientation action:", err)
		return action, err
	}

	return action, nil

}
