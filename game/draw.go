package game

import (
	"fmt"
	"image" // Ajouter par moi
	"image/color"
	"log"
	"math"
	"sort"
	"strconv"
	"time" // Ajouter par moi pour le cycle jour/nuit

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/character"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/assets" // AJouter par moi
)

type charArr []character.Character

func (cArr charArr) Len() int {
	return len(cArr)
}
func (cArr charArr) Swap(i, j int) {
	cArr[i], cArr[j] = cArr[j], cArr[i]
}
func (cArr charArr) Less(i, j int) bool {
	return !cArr[i].IsInFrontOf(&cArr[j])
}

// Draw permet d'afficher à l'écran tous les éléments du jeu
// (le sol, le personnage, les éventuelles informations de debug).
// Il faut faire attention à l'ordre d'affichage pour éviter d'avoir
// des éléments qui en cachent d'autres.

func (g *Game) Draw(screen *ebiten.Image) {
	// Daycycle (cycle jour/nuit)
	if configuration.Global.Daycycle {

		//Etape1 : récupérer l'heure actuelle
		currentTime := time.Now()
		heures := float64(currentTime.Hour())
		minutes := float64(currentTime.Minute())
		secondes := float64(currentTime.Second())

		//Etape2 : Faire des calculs pour que toutes les 20 minutes, on passe du jour, à la nuit
		temps := (heures*60 + minutes + secondes/60)
		brouillon := (temps - float64(int(temps)/40)*40) - 20
		if brouillon < 0 {
			brouillon = brouillon * -1
		}
		new_br := ((1.0 - float64(brouillon)/40) * 2) - 1.5

		if new_br > 0 {
			new_br = 0
		}
		configuration.Global.AmbiantLight = new_br

	}

	// Déplacement par molette dans la hotbar

	if configuration.Global.HotbarOn && g.PlayerOnUi == "" {
		// Les test fait au dessus permettent de vérifier que le joueur ne se trouve pas sur une interface ainsi que l'inventaire à bien été
		// activer dans le fichier configuration.json

		_, wheelY := ebiten.Wheel() // Récupérer les mouvements de la molette

		// On a récupérer la position de la molette pour pouvoir changer de slot dans la hotbar par son utilisation
		// PS : Fonctionne très bien mais la vitesse de changement de slot dépend de la souris et de l'utilisation d'un pad
		if wheelY < 0 {
			if configuration.Global.HotbarSelect > 7 {
				configuration.Global.HotbarSelect = 0
			} else {
				configuration.Global.HotbarSelect += 1
			}
		}
		if wheelY > 0 {
			if configuration.Global.HotbarSelect < 1 {
				configuration.Global.HotbarSelect = 8
			} else {
				configuration.Global.HotbarSelect -= 1
			}
		}
	}

	// Charger les particules de four pour facilité la tache
	Particle, _, _ := ebitenutil.NewImageFromFile("../assets/particles/no_particle.png") // Si jamais on modifie le json pendant, ça plantera pas ici
	configuration.Global.FireParticle = false
	if configuration.Global.ParticleOn {
		Particle, _, _ = ebitenutil.NewImageFromFile("../assets/particles/fire.png")
		configuration.Global.FireParticle = true
	}

	// Ne pas faire attention, les deux lignes d code récupèrent des informations pour la création future d'une caméra fluide
	x_add, y_add := g.Character.GetContentPLayerAddToCamPos()
	cam_x, cam_y := g.Camera.Get_move_bool_allowed()

	PlaceBlocX, PlaceBlocY, PlaceBloc, PlaceBlocID := g.Floor.Draw(screen, Particle, g.Character.X, g.Character.Y, g.Camera.Y, g.Camera.X, x_add, y_add, cam_x, cam_y, g.PlayerOnUi)

	if PlaceBloc {
		g.PlaceBlocX, g.PlaceBlocY, g.PlaceBloc, g.PlaceBlocID = PlaceBlocX, PlaceBlocY, PlaceBloc, PlaceBlocID
	}

	// Multijoueur : crée un array pour stocker les joueurs
	allChar := make([]character.Character, len(g.OtherPlayers)+1)
	allChar[0] = g.Character
	for i := range g.OtherPlayers {
		allChar[i+1] = g.OtherPlayers[i]
	}
	sort.Sort(charArr(allChar)) // S'assurer que les joueurs soient triés pour apparaître devant où derrière selon leur index dans l'array

	// Multijoueur : Dessiner tous les joueurs sur la map
	for _, player := range allChar {
		if player.Name == configuration.Global.PlayerName {
			player.Draw(screen, g.Camera.X, g.Camera.Y, cam_x, cam_y)
		} else {
			player.Draw(screen, g.Camera.X, g.Camera.Y, true, true)
		}
	}

	if !configuration.Global.DebugMode { // Pour ne pas déranger la vue en mode de débeuguage, on cache les interfaces
		// Ici on s'occupe d'afficher la majorité des différentes ui pour le joueur

		//Ajout de la hotbar
		if configuration.Global.HotbarOn {
			drawAt("../assets/hotbar.png", 50, 90, 90, screen)

			hotbarXPos := 10 + float64(configuration.Global.HotbarSelect)*10
			drawAt("../assets/hotbar_select.png", hotbarXPos, 90, 11, screen)

			//Items dans la hotbar
			for i := 0; i < 9; i++ {
				hotbarXPos := 10 + float64(i)*10
				drawAtBUTbloc(configuration.Global.InventoryItems[i][0], hotbarXPos, 90, 7, screen)
			}
		}

		//Ajout de la Hungerbar
		if configuration.Global.HungerbarOn {
			imagePath := "../assets/hungerbar/" + strconv.Itoa(configuration.Global.HungerbarLevel) + ".png"
			drawAt(imagePath, 75, 80, 30, screen)
		}

		//Affichage de tous ce qui est ui
		if g.PlayerOnUi != "" {

			// Assombrir l'arrière plan
			drawAt("../assets/UI/black.png", 50, 50, 100, screen)

			//Recherche du type de ui à afficher

			if g.PlayerOnUi == "inventory" && configuration.Global.UiInventaireOn {
				//Affichage de l'inventaire
				drawAt("../assets/UI/inventory.png", 50, 50, 70, screen)
			}
		}
	}

	//mode de débeugage (s'active avec F3)
	if configuration.Global.DebugMode {
		g.drawDebug(screen)
	}
}

// drawDebug se charge d'afficher les informations de debug si
// l'utilisateur le demande (positions absolues du personnage
// et de la caméra, grille avec les coordonnées, etc).
func (g Game) drawDebug(screen *ebiten.Image) {
	// Cette fonction s'occupe de l'affichage du débeuguage

	gridColor := color.NRGBA{R: 255, G: 255, B: 255, A: 63}
	gridHoverColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	gridLineSize := 2
	cameraColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	cameraLineSize := 1

	mouseX, mouseY := ebiten.CursorPosition()

	xMaxPos := configuration.Global.ScreenWidth * 4
	yMaxPos := configuration.Global.ScreenHeight * 4

	//Dessine les lignes verticales
	for x := 0; x < configuration.Global.RealNumTileX; x++ {
		xGeneralPos := x*configuration.Global.TileSize*4 + configuration.Global.TileSize*4/2
		xPos := float32(xGeneralPos)

		lineColor := gridColor
		if mouseX+1 >= xGeneralPos && mouseX+1 <= xGeneralPos+gridLineSize {
			lineColor = gridHoverColor
		}

		vector.StrokeLine(screen, xPos, 0, xPos, float32(yMaxPos), float32(gridLineSize), lineColor, false)

		xPrintValue := g.Camera.X + x - configuration.Global.ScreenCenterTileX
		xPrint := fmt.Sprint(xPrintValue)
		if len(xPrint) <= (2*configuration.Global.TileSize)/16 || (xPrintValue > 0 && xPrintValue%2 == 0) || (xPrintValue < 0 && (-xPrintValue)%2 == 0) {
			xTextPos := xGeneralPos - 3*len(xPrint) - 1
			ebitenutil.DebugPrintAt(screen, xPrint, xTextPos, yMaxPos)
		}
	}
	//Dessine les lignes horizontales
	for y := 0; y < configuration.Global.RealNumTileY; y++ {
		yGeneralPos := y*configuration.Global.TileSize + configuration.Global.TileSize/2
		yPos := float32(yGeneralPos) * 4

		lineColor := gridColor
		if mouseY+1 >= yGeneralPos && mouseY+1 <= yGeneralPos+gridLineSize {
			lineColor = gridHoverColor
		}

		vector.StrokeLine(screen, 0, yPos, float32(xMaxPos), yPos, float32(gridLineSize), lineColor, false)

		yPrint := fmt.Sprint(g.Camera.Y + y - configuration.Global.ScreenCenterTileY)
		xTextPos := xMaxPos + 1
		yTextPos := yGeneralPos - 8
		ebitenutil.DebugPrintAt(screen, yPrint, xTextPos, yTextPos)
	}
	// Fait le rectangle pour la caméra en debug mode
	vector.StrokeRect(screen, float32(configuration.Global.ScreenCenterTileX*configuration.Global.TileSize*4), float32(configuration.Global.ScreenCenterTileY*configuration.Global.TileSize*4), float32((configuration.Global.TileSize*4)+1), float32((configuration.Global.TileSize*4)+1), float32(cameraLineSize), cameraColor, false)

	// Ecrit le texte du debug mode

	//PS: J'ai préférer ne pas trop toucher en dessous, mais c'est très dérangeant d'avoir les coo x et y inversées
	ebitenutil.DebugPrintAt(screen, "Camera:", 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprint("(", g.Camera.X, ",", g.Camera.Y, ")"), 0, int(float64(configuration.Global.ScreenWidth)*0.1))

	ebitenutil.DebugPrintAt(screen, "Character:", 0, int(float64(configuration.Global.ScreenWidth)*0.2))
	ebitenutil.DebugPrintAt(screen, fmt.Sprint("(", g.Character.X, ",", g.Character.Y, ")"), 0, int(float64(configuration.Global.ScreenWidth)*0.3))

}

func drawAt(imagePath string, xPourcentage, yPourcentage, sizePourcentage float64, screen *ebiten.Image) {
	// Cette fonction sert à charger une image et à l'afficher peut importe les différents paramètres choisis
	// Elle permet au passage d'imaginer un ancrage au milieu de l'image même si en vérité, l'ancrage d'une image
	// se trouve en haut à gauche

	// Etape 1 : Chargement de l'image
	Photbar, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		log.Fatal("Erreur lors du chargement de l'image:", err)
		return
	}

	// Etape 2 : Calculs de la taille et de la position par rapport aux pourcentages reçus et la taille réelle de la fenêtre et de la map
	imageWidth, imageHeight := Photbar.Size()  // Taille de l'image
	screenWidth, screenHeight := screen.Size() // Taille de l'écran

	targetWidth := float64(screenWidth) * sizePourcentage / 100
	targetHeight := float64(screenHeight) * sizePourcentage / 100
	uniformScale := math.Min(targetWidth/float64(imageWidth), targetHeight/float64(imageHeight))
	scaledWidth := float64(imageWidth) * uniformScale
	scaledHeight := float64(imageHeight) * uniformScale

	targetX := xPourcentage * float64(screenWidth) / 100
	targetY := yPourcentage * float64(screenHeight) / 100
	finalX := targetX - scaledWidth/2
	finalY := targetY - scaledHeight/2

	// Configurer les options puis afficher l'image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(uniformScale, uniformScale)
	op.GeoM.Translate(finalX, finalY)

	screen.DrawImage(Photbar, op)
}

func drawAtBUTbloc(blocID int, xPourcentage, yPourcentage, sizePourcentage float64, screen *ebiten.Image) {
	// Cette fonction sert à charger une image (à partir du blocID donné en paramètre #le bloc id est un numéro attribué à chaque bloc selon sa position)
	// et à l'afficher peut importe les différents paramètres choisis.
	// Elle permet au passage d'imaginer un ancrage au milieu de l'image même si en vérité, l'ancrage d'une image
	// se trouve en haut à gauche

	// PS: C'est à peu près la même fonction que celle juste au dessus, mais celle-ci s'occupe seulement des blocs

	// Etape 1 : Calculs de la taille et de la position par rapport aux pourcentages reçus et la taille réelle de la fenêtre et de la map
	imageWidth, imageHeight := configuration.Global.TileSize, configuration.Global.TileSize // Taille d'un bloc
	screenWidth, screenHeight := screen.Size()                                              // Taille de l'écran

	targetWidth := float64(screenWidth) * sizePourcentage / 100
	targetHeight := float64(screenHeight) * sizePourcentage / 100
	uniformScale := math.Min(targetWidth/float64(imageWidth), targetHeight/float64(imageHeight))
	scaledWidth := float64(imageWidth) * uniformScale
	scaledHeight := float64(imageHeight) * uniformScale

	targetX := xPourcentage * float64(screenWidth) / 100
	targetY := yPourcentage * float64(screenHeight) / 100
	finalX := targetX - scaledWidth/2
	finalY := targetY - scaledHeight/2

	// Etape 2 : Configurer les options puis afficher l'image
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(uniformScale, uniformScale)
	op.GeoM.Translate(finalX, finalY)

	shiftX := blocID * configuration.Global.TileSize
	screen.DrawImage(assets.FloorImage.SubImage(
		image.Rect(shiftX, 0, shiftX+configuration.Global.TileSize, configuration.Global.TileSize),
	).(*ebiten.Image), op)
}
