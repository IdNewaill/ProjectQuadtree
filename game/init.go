package game

import (
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
)

// Init initialise les données d'un jeu. Il faut bien
// faire attention à l'ordre des initialisation car elles
// pourraient dépendre les unes des autres.
func (g *Game) Init(seed int64, playerX, playerY int, playerName string) {
	g.Floor.Init(configuration.Global.Seed)
	g.Character.Init(playerX, playerY, playerName)
	g.Camera.Init(playerX, playerY)

	// Variables de fonctionnement
	g.PlayerOnUi = ""
}
