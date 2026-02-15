package character

// Init met en place un personnage. Pour le moment
// cela consiste simplement à initialiser une variable
// responsable de définir l'étape d'animation courante.
func (c *Character) Init(x, y int, playerName string) {
	c.animationStep = 1

	c.X = x
	c.Y = y
	c.Name = playerName
}
