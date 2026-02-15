package game

import (
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_out"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/floor"
)

type ChunkLoaded struct {
	coord floor.Coords
	chunk floor.Chunk
}

func LoadChunk(x, y int, sendTo string, actionOutChan chan ActionOut, chunkLoadedChan chan ChunkLoaded) {
	// Cette fonction permet de charger un chunk et les infos nécessaires

	content := floor.LoadChunkContent(x, y, configuration.Global.GameDir)

	if sendTo != "" {
		actionOutChan <- action_out.NewSetChunk(sendTo, x, y, content)
	}

	chunkLoadedChan <- ChunkLoaded{
		coord: floor.Coords{X: x, Y: y},
		chunk: floor.NewChunk(content),
	}
}
