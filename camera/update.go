package camera

import (
	"bufio"
	"os"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
)

// Update met à jour la position de la caméra à chaque pas
// de temps, c'est-à-dire tous les 1/60 secondes.
func (c *Camera) Update(characterPosX, characterPosY int, x_add, y_add int) {

	switch configuration.Global.CameraMode {
	case Static:
		c.updateStatic()
	case FollowCharacter:
		c.updateFollowCharacter(characterPosX, characterPosY, x_add, y_add)
	case LimiteTerrain:
		c.updateLimiteTerrain(characterPosX, characterPosY, x_add, y_add)
	}
}

// updateStatic est la mise-à-jour d'une caméra qui reste
// toujours à la position (0,0). Cette fonction ne fait donc
// rien.
func (c *Camera) updateStatic() {}

// updateFollowCharacter est la mise-à-jour d'une caméra qui
// suit toujours le personnage. Elle prend en paramètres deux
// entiers qui indiquent les coordonnées du personnage et place
// la caméra au même endroit.
func (c *Camera) updateFollowCharacter(characterPosX, characterPosY, x_add, y_add int) {
	c.X = characterPosX
	c.Y = characterPosY
}

var x_move_bool bool
var y_move_bool bool

// updateFollowCharacter met à jour la position de la caméra en suivant le personnage,
// tout en le bloquant aux bords du terrain.(Fait par Shawn modifié par Grégoire)
func (c *Camera) updateLimiteTerrain(characterPosX, characterPosY, x_add, y_add int) {

	readFile, err := os.Open(configuration.Global.FloorFile)
	if err != nil {
		os.Exit(0)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	var taille_y int = 0
	var taille_x int = 0
	for fileScanner.Scan() {
		taille_y += 1
		line := fileScanner.Text()
		taille_x = len(line)
	}
	taille_y -= 1
	taille_x -= 1

	screenWidth := configuration.Global.NumTileX
	screenHeight := configuration.Global.NumTileY
	terrainWidth := taille_x / configuration.Global.SizeBlockID
	terrainHeight := taille_y

	if screenWidth > terrainWidth || screenHeight > terrainHeight {
		if terrainHeight > terrainWidth {
			configuration.Global.NumTileX, configuration.Global.NumTileY = terrainWidth, terrainWidth
		} else {
			configuration.Global.NumTileX, configuration.Global.NumTileY = terrainHeight, terrainHeight
		}
	}

	screenWidth = configuration.Global.NumTileX
	screenHeight = configuration.Global.NumTileY
	halfScreenWidth := screenWidth / 2
	halfScreenHeight := screenHeight / 2

	c.X = characterPosX
	c.Y = characterPosY

	if c.X < halfScreenWidth {
		c.X = halfScreenWidth
		x_move_bool = true
	} else if c.X > terrainWidth-halfScreenWidth {
		c.X = terrainWidth - halfScreenWidth
		x_move_bool = true
	} else {
		x_move_bool = false
	}

	if c.Y < halfScreenHeight {
		c.Y = halfScreenHeight
		y_move_bool = true
	} else if c.Y > terrainHeight-halfScreenHeight {
		c.Y = terrainHeight - halfScreenHeight
		y_move_bool = true
	} else {
		y_move_bool = false
	}
}

func (c *Camera) Get_move_bool_allowed() (bool, bool) {
	return x_move_bool, y_move_bool
}
