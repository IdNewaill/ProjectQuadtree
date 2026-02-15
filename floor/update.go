package floor

import (
	"fmt"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
)

// Update se charge de stocker dans la structure interne (un tableau)
// de f une représentation de la partie visible du terrain à partir
// des coordonnées absolues de la case sur laquelle se situe la
// caméra.

// Taux de refresh
var ancienne_valeur_secondes int64

func (f *Floor) Update(camXPos, camYPos int) []Coords {

	for y := 0; y < configuration.Global.RealNumTileY; y++ {
		for x := 0; x < configuration.Global.RealNumTileX; x++ {
			f.content[y][x] = -1
		}
	}
	topLeftX := camXPos - configuration.Global.ScreenCenterTileX
	topLeftY := camYPos - configuration.Global.ScreenCenterTileY
	bottomRightX := topLeftX + configuration.Global.NumTileX
	bottomRightY := topLeftY + configuration.Global.NumTileY

	if configuration.Global.MapFileMode {
		switch configuration.Global.FloorKind {
		case 0:
			f.updateGridFloor(topLeftX, topLeftY)
		case 1:
			f.updateFromFileFloor(topLeftX, topLeftY)
		case 2:
			f.updateQuadtreeFloor(topLeftX, topLeftY)
		}
		return nil
	} else {
		chunkStartX := topLeftX / 64
		if topLeftX < 0 {
			chunkStartX -= 1
		}

		chunkEndX := bottomRightX / 64
		if bottomRightX < 0 {
			chunkEndX -= 1
		}

		chunkStartY := topLeftY / 64
		if topLeftY < 0 {
			chunkStartY -= 1
		}

		chunkEndY := bottomRightY / 64
		if bottomRightY < 0 {
			chunkEndY -= 1
		}

		chunkToLoad := []Coords{}
		if ReloadChunk {
			chunkToLoad = append(chunkToLoad, Coords{X: ReloadChunkX, Y: ReloadChunkY})
			ReloadChunk = false
		}
		for x := chunkStartX; x <= chunkEndX; x++ {
			for y := chunkStartY; y <= chunkEndY; y++ {

				chunk, ok := f.fullContent[Coords{x, y}]
				if !ok {

					// regarde si le chunk n'a pas déjà été demander
					isChunkPending := false
					for k := range f.pendingChunks {
						if f.pendingChunks[k].X == x && f.pendingChunks[k].Y == y {
							isChunkPending = true
							break
						}
					}
					if !isChunkPending {
						coord := Coords{X: x, Y: y}
						chunkToLoad = append(chunkToLoad, coord)
						f.pendingChunks = append(f.pendingChunks, coord)
					}

					continue
				}
				/*  // Sert pour le multijouer quand on pose des blocs mais fonctionne très mal
				if configuration.Global.GameType == "join" {
					now := time.Now()
					secondes := now.Unix()
					if secondes > ancienne_valeur_secondes {
						coord := Coords{X: x, Y: y}
						chunkToLoad = append(chunkToLoad, coord)
						f.pendingChunks = append(f.pendingChunks, coord)
						ancienne_valeur_secondes = secondes
					}
				}*/

				chunk.quadtree.GetContent(
					topLeftX-x*64,
					topLeftY-y*64,
					f.content,
				)
			}
		}
		if len(chunkToLoad) > 0 {
			fmt.Println(chunkToLoad)
		}
		return chunkToLoad
	}
}

// ################ Anciennes fonctions ici ##############

// le sol est un quadrillage de tuiles d'herbe et de tuiles de désert
func (f *Floor) updateGridFloor(topLeftX, topLeftY int) {
	for y := 0; y < len(f.content); y++ {
		for x := 0; x < len(f.content[y]); x++ {
			absX := topLeftX
			if absX < 0 {
				absX = -absX
			}
			absY := topLeftY
			if absY < 0 {
				absY = -absY
			}
			f.content[y][x] = ((x + absX%2) + (y + absY%2)) % 2
		}
	}
}

// le sol est récupéré depuis un tableau, qui a été lu dans un fichier
//
// la version actuelle recopie fullContent dans content, ce qui n'est pas
// le comportement attendu dans le rendu du projet
func (f *Floor) updateFromFileFloor(topLeftX, topLeftY int) {
	for y := 0; y < configuration.Global.RealNumTileY; y++ {
		for x := 0; x < configuration.Global.RealNumTileX; x++ {
			realY := y + topLeftY
			realX := x + topLeftX

			if realY >= 0 && realY < len(f.NORMALfullContent) && realX >= 0 && realX < len(f.NORMALfullContent[realY]) {
				f.content[y][x] = f.NORMALfullContent[realY][realX]
			} else {
				f.content[y][x] = -1 // remplit avec -1 si hors des limites du terrain
			}
		}
	}
}

// le sol est récupéré depuis un quadtree, qui a été lu dans un fichier
func (f *Floor) updateQuadtreeFloor(topLeftX, topLeftY int) {
	f.quadtreeContent.GetContent(topLeftX, topLeftY, f.content)
}
