package game_action

import (
	"fmt"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
)

type DespawnPlayer struct {
	Name string
}

func (a DespawnPlayer) Execute(g *game.Game) {

	for index := range g.OtherPlayers {
		if g.OtherPlayers[index].Name == a.Name {

			fmt.Println("player", a.Name, "disconnected")
			g.OtherPlayers[index] = g.OtherPlayers[len(g.OtherPlayers) - 1]
			g.OtherPlayers = g.OtherPlayers[:len(g.OtherPlayers)-1]

			return
		}
	}
}

func decodeDespawnPlayer(data []byte) (DespawnPlayer, error) {
	return DespawnPlayer{
		Name: string(data),
	}, nil
}