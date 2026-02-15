package game_action

import (
	"errors"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
)

func GetRemoteActionFrom(data []byte) (game.GameAction, error) {
	if len(data) == 0 {
		return nil, errors.New("cannot decode empty data")
	}

	id := data[0]

	switch id {
	case action_type.PlayerSpawnId:
		return decodeSpawnPlayer(data[1:])

	case action_type.PlayerDespawnId:
		return decodeDespawnPlayer(data[1:])

	case action_type.PlayerMoveId:
		return decodePlayerMove(data[1:])

	case action_type.OrientationChangeId:
		return decodeOrientationChange(data[1:])

	case action_type.AskChunkActionId:
		return decodeAskChunk(data[1:])

	case action_type.SetChunkActionId:
		return decodeSetChunk(data[1:])

	default:
		return nil, errors.New("unknown id")
	}
}