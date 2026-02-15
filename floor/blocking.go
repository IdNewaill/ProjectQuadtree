package floor

import (
	"fmt"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
)

// Blocking retourne, étant donnée la position du personnage,
// un tableau de booléen indiquant si les cases au dessus (0),
// à droite (1), au dessous (2) et à gauche (3) du personnage
// sont bloquantes.
func (f Floor) Blocking(characterXPos, characterYPos, camXPos, camYPos int) (blocking [6]bool) {
	relativeXPos := characterXPos - camXPos + configuration.Global.ScreenCenterTileX
	relativeYPos := characterYPos - camYPos + configuration.Global.ScreenCenterTileY

	Ice := configuration.Global.IceID
	Cobweb := configuration.Global.CobwebID
	Stair := configuration.Global.StairID
	StillClimbing := false

	if configuration.Global.NewBlocking {
		for _, v := range configuration.Global.Blocking {
			if v == f.content[relativeYPos][relativeXPos] {
				StillClimbing = true
			}
		}
		if !StillClimbing {
			configuration.Global.OnStair = false
		}
	}

	if f.content[relativeYPos][relativeXPos] == Stair {
		configuration.Global.OnStair = true
	}

	// Calcul des cases bloquantes, glissante ou collant
	if relativeXPos > 0 && relativeYPos > 0 {
		if configuration.Global.CanClimb && (StillClimbing || configuration.Global.OnStair) {
			fmt.Println(relativeXPos, relativeYPos)
			blocking[0] = f.content[relativeYPos-1][relativeXPos] == -1
			blocking[1] = f.content[relativeYPos][relativeXPos+1] == -1
			blocking[2] = f.content[relativeYPos+1][relativeXPos] == -1
			blocking[3] = f.content[relativeYPos][relativeXPos-1] == -1
		} else {
			blocking[0] = CanWalk(relativeYPos-1, relativeXPos, f)
			blocking[1] = CanWalk(relativeYPos, relativeXPos+1, f)
			blocking[2] = CanWalk(relativeYPos+1, relativeXPos, f)
			blocking[3] = CanWalk(relativeYPos, relativeXPos-1, f)
		}
	}
	if configuration.Global.Slide {
		blocking[4] = f.content[relativeYPos][relativeXPos] == Ice
	}

	if configuration.Global.AffectedByCobweb {
		blocking[5] = f.content[relativeYPos][relativeXPos] == Cobweb
	}

	return blocking
}

// Vérifie que le character sera sur le terrain puis utilise la fonction contains
func CanWalk(characterFuturYPos, characterFuturXPos int, f Floor) bool {
	if characterFuturXPos >= 0 && characterFuturXPos < len(f.content[0]) && characterFuturYPos >= 0 && characterFuturYPos < len(f.content) {
		return contains(f.content[characterFuturYPos][characterFuturXPos])
	}
	return true
}

// Vérifie que le bloc de la future position n'est pas bloquant
func contains(value int) bool {
	slice := configuration.Global.Blocking
	if configuration.Global.NewBlocking {
		for _, v := range slice {
			if v == value {
				return true
			}
		}
	} else if value == -1 {
		return true
	}
	return false
}
