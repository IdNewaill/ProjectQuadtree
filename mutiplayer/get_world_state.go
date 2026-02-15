package game_action

import (
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
)

type GetWorldState struct {
	Sender chan action_type.InitWorldState
}

func (a GetWorldState) Execute(g *game.Game) {

	players := make([]action_type.PlayerInitData, len(g.OtherPlayers) + 1)
	players[0] = action_type.PlayerInitData{
		Name: g.Character.Name,
		X: g.Character.X,
		Y: g.Character.Y,
		Orientation: g.Character.Orientation,
	}

	for i := range g.OtherPlayers {
		players[i + 1] = action_type.PlayerInitData{
			Name: g.OtherPlayers[i].Name,
			X: g.OtherPlayers[i].X,
			Y: g.OtherPlayers[i].Y,
			Orientation: g.OtherPlayers[i].Orientation,
		}
	}
	
	a.Sender <- action_type.InitWorldState {
		Players: players,
	}
}