package camera

// Init initialise la caméra avec les coordonnées du personnage et les dimensions de l'écran et du terrain.
func (c *Camera) Init(characterPosX, characterPosY int) {
	c.X = characterPosX
	c.Y = characterPosY
}
