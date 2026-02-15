package character

import (
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"

	"github.com/hajimehoshi/ebiten/v2"
)

// Update met à jour la position du personnage, son orientation
// et son étape d'animation (si nécessaire) à chaque pas
// de temps, c'est-à-dire tous les 1/60 secondes.

func (c *Character) Update(blocking [6]bool) (moving bool) {
	if !c.moving {
		if (ebiten.IsKeyPressed(ebiten.KeyRight) && configuration.Global.OriginalKey) || (ebiten.IsKeyPressed(ebiten.KeyD) && !configuration.Global.OriginalKey) {
			c.Orientation = orientedRight
			if !blocking[1] {
				c.xInc = 1
				c.moving = true
			}
		} else if (ebiten.IsKeyPressed(ebiten.KeyLeft) && configuration.Global.OriginalKey) || (ebiten.IsKeyPressed(ebiten.KeyQ) && !configuration.Global.OriginalKey) || (ebiten.IsKeyPressed(ebiten.KeyA) && !configuration.Global.OriginalKey) {
			c.Orientation = orientedLeft
			if !blocking[3] {
				c.xInc = -1
				c.moving = true

			}
		} else if (ebiten.IsKeyPressed(ebiten.KeyUp) && configuration.Global.OriginalKey) || (ebiten.IsKeyPressed(ebiten.KeyZ) && !configuration.Global.OriginalKey) || (ebiten.IsKeyPressed(ebiten.KeyW) && !configuration.Global.OriginalKey) {
			c.Orientation = orientedUp
			if !blocking[0] {
				c.yInc = -1
				c.moving = true
			}
		} else if (ebiten.IsKeyPressed(ebiten.KeyDown) && configuration.Global.OriginalKey) || (ebiten.IsKeyPressed(ebiten.KeyS) && !configuration.Global.OriginalKey) {
			c.Orientation = orientedDown
			if !blocking[2] {
				c.yInc = 1
				c.moving = true
			}
		}
	} else {
		//Change la vitesse de l'animation pour différents types de blocs
		var SpeedCoefficient float32 = 1
		if blocking[5] && configuration.Global.AffectedByCobweb {
			SpeedCoefficient = 6
		}
		if blocking[4] && configuration.Global.Slide {
			SpeedCoefficient = 0.5
		}
		c.animationFrameCount++
		if c.animationFrameCount >= int(float32(configuration.Global.NumFramePerCharacterAnimImage)*SpeedCoefficient) {
			c.animationFrameCount = 0
			shiftStep := configuration.Global.TileSize / configuration.Global.NumCharacterAnimImages
			c.shift += shiftStep
			c.animationStep = -c.animationStep
			if c.shift > configuration.Global.TileSize-shiftStep {
				c.shift = 0
				c.moving = false
				c.X += c.xInc
				c.Y += c.yInc
				c.xInc = 0
				c.yInc = 0
			}
		}
	}
	return c.moving
}
