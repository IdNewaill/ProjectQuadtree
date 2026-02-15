package character

import "fmt"

const (
	orientedDown int = iota
	orientedLeft
	orientedRight
	orientedUp
)

// Character définit les charactéristiques du personnage.
// Pour le moment seules les coordonnées absolues de l'endroit
// où il se trouve sont exportées, mais vous pourrez
// ajouter des choses au besoins lors de votre développement.
//
// Les champs non exportés sont :
//   - orientation : l'orientation du personnage (haut, bas, gauche, droite).
//   - animationStep : l'étape d'animation (-1 ou 1, représentant l'animation
//     d'un pas à gauche ou à droite).
//   - xInc, yInc : les incréments en X et Y à réaliser après la prochaine animation.
//   - moving : l'information de si une animation est en cours ou pas.
//   - shift : la position actuelle en pixels du personnage relativement à ses
//     coordonnées absolues.
//   - animationFrameCount : le nombre d'appels à update (ou de 1/60 de seconde) qui
//     ont eu lieu depuis la dernière étape d'animation.
type Character struct {
	Name                string
	X, Y                int
	Orientation         int
	animationStep       int
	xInc, yInc          int
	moving              bool
	shift               int
	animationFrameCount int
}

// Détermine si le character est devant un autre joueur en multijoueur par exemple.
func (c0 *Character) IsInFrontOf(c1 *Character) bool {
	if c0.Y != c1.Y {
		return c0.Y > c1.Y
	}

	if !c0.moving && c1.moving {
		return c1.Orientation == orientedUp
	}

	if c0.moving && !c1.moving {
		return c0.Orientation == orientedDown
	}

	if c0.moving && c1.moving {
		if c0.Orientation == orientedDown && c1.Orientation == orientedDown {
			return c0.shift > c1.shift
		}

		if c0.Orientation == orientedUp && c1.Orientation == orientedUp {
			return c0.shift < c1.shift
		}
	}

	return true
}

// Regarde si le character peut aller dans la direction qu'il veut aller, s'il ne peut pas, il change son orientation
func (c *Character) TryMoveUp(blocking [6]bool) bool {
	if c.moving {
		return false
	}

	c.Orientation = orientedUp

	if blocking[0] {
		return false
	}

	c.yInc = -1
	c.moving = true
	return true
}
func (c *Character) TryMoveRight(blocking [6]bool) bool {
	if c.moving {
		return false
	}

	c.Orientation = orientedRight

	if blocking[1] {
		return false
	}

	c.xInc = 1
	c.moving = true
	return true
}

func (c *Character) TryMoveDown(blocking [6]bool) bool {
	if c.moving {
		return false
	}

	c.Orientation = orientedDown

	if blocking[2] {
		return false
	}

	c.yInc = 1
	c.moving = true
	return true
}

func (c *Character) TryMoveLeft(blocking [6]bool) bool {
	if c.moving {
		return false
	}

	c.Orientation = orientedLeft

	if blocking[3] {
		return false
	}

	c.xInc = -1
	c.moving = true
	return true
}

// Mets à jour la position des autres joueurs en multijoueur
/*func (c *Character) RemoteMove(newX, newY int) {
	NO_IA
	switch {
	case newX > c.X:
		c.X = newX - 1
		c.Orientation = orientedRight
		c.xInc = 1
	case newX < c.X:
		c.X = newX + 1
		c.Orientation = orientedLeft
		c.xInc = -1
	case newY > c.Y:
		c.Y = newY - 1
		c.Orientation = orientedDown
		c.yInc = 1
	case newY < c.Y:
		c.Y = newY + 1
		c.Orientation = orientedUp
		c.yInc = -1
	}

	c.moving = true
	c.shift = 0

}*/

func (c *Character) RemoteMove(newX, newY int) {
	c.shift = 0
	c.moving = true
	fmt.Println("Bouger")
	if newX > c.X {
		c.Orientation = orientedRight
		c.xInc = 1
		c.yInc = 0
	} else if newX < c.X {
		c.Orientation = orientedLeft
		c.xInc = -1
		c.yInc = 0
	} else if newY > c.Y {
		c.Orientation = orientedDown
		c.yInc = 1
		c.xInc = 0
	} else if newY < c.Y {
		c.Orientation = orientedUp
		c.yInc = -1
		c.xInc = 0
	}
}

// Mets à jour l'orientation des autres joueurs en multijoueur
func (c *Character) RemoteSetOrientation(newOrientation int) {
	c.Orientation = newOrientation
}

// Créer un nouveau joueur dans le jeu local pour le multijoueur
func New(name string, x, y int) Character {
	return Character{
		Name:          name,
		animationStep: 1,
		X:             x,
		Y:             y,
	}
}
