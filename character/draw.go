package character

import (
	"image"
	"image/color"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/assets"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Draw permet d'afficher le personnage dans une *ebiten.Image
// (en pratique, celle qui représente la fenêtre de jeu) en
// fonction des charactéristiques du personnage (position, orientation,
// étape d'animation, etc) et de la position de la caméra (le personnage
// est affiché relativement à la caméra).

var KeepxTileForDisplay int
var KeepyTileForDisplay int

func (c Character) Draw(screen *ebiten.Image, camX, camY int, can_x, can_y bool) {

	xShift := 0
	yShift := 0
	var PlayerOrientation int
	switch c.Orientation {
	case orientedDown:
		yShift = c.shift
		PlayerOrientation = 2
	case orientedUp:
		yShift = -c.shift
		PlayerOrientation = 1
	case orientedLeft:
		xShift = -c.shift
		PlayerOrientation = 0
	case orientedRight:
		xShift = c.shift
		PlayerOrientation = 3
	}

	xTileForDisplay := c.X - camX + configuration.Global.ScreenCenterTileX
	yTileForDisplay := c.Y - camY + configuration.Global.ScreenCenterTileY

	xPos := (xTileForDisplay)*configuration.Global.TileSize + xShift
	yPos := (yTileForDisplay)*configuration.Global.TileSize - configuration.Global.TileSize/2 + 2 + yShift
	if false && xPos-(xPos/16)*16 < 12 && yPos+6-((yPos+6)/16)*16 < 12 {

		if PlayerOrientation != 0 {
			KeepxTileForDisplay = xPos - (xPos/16)*16
			xPos -= KeepxTileForDisplay
		} else {
			KeepxTileForDisplay = (xPos - 3 - ((xPos-3)/16)*16 - 16) / 2
			xPos -= KeepxTileForDisplay
		}

		if PlayerOrientation == 1 {
			KeepyTileForDisplay = (yPos + 3 - ((yPos+3)/16)*16 - 16) / 2
			yPos -= KeepyTileForDisplay
		} else {
			KeepyTileForDisplay = yPos + 6 - ((yPos+6)/16)*16
			yPos -= KeepyTileForDisplay
		}
	}
	//xPos -= KeepxTileForDisplay
	//yPos -= KeepyTileForDisplay
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xPos), float64(yPos))

	shiftX := configuration.Global.TileSize
	if c.moving {
		shiftX += c.animationStep * configuration.Global.TileSize
	}
	shiftY := c.Orientation * configuration.Global.TileSize
	op.GeoM.Scale(4, 4)

	screen.DrawImage(assets.CharacterImage.SubImage(
		image.Rect(shiftX, shiftY, shiftX+configuration.Global.TileSize, shiftY+configuration.Global.TileSize),
	).(*ebiten.Image), op)

	// draw name
	nameBound := text.BoundString(font.Minecraft, c.Name)
	nameX := (xPos*4 + configuration.Global.TileSize*2) - nameBound.Dx()/2
	nameY := yPos*4 - 10

	vector.DrawFilledRect(
		screen,
		float32(nameX-4),
		float32(nameY-nameBound.Dy()-4),
		float32(nameBound.Dx()+8),
		float32(nameBound.Dy()+8),
		color.RGBA{R: 0, G: 0, B: 0, A: 128},
		false,
	)

	text.Draw(
		screen,
		c.Name,
		font.Minecraft,
		nameX,
		nameY,
		color.White,
	)

}

func (c Character) GetContentCharacterPos(camX, camY int) (x, y int) { //Ajoutée par moi

	xTileForDisplay := c.X - camX + configuration.Global.ScreenCenterTileX
	yTileForDisplay := c.Y - camY + configuration.Global.ScreenCenterTileY
	return xTileForDisplay, yTileForDisplay
}

func (c Character) GetContentPLayerAddToCamPos() (add_x, add_y int) {
	return KeepxTileForDisplay, KeepyTileForDisplay
}
