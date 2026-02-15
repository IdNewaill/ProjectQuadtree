package floor

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/assets"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"

	"math/rand" //Ajouter par moi pour le random
	//"github.com/hajimehoshi/ebiten/v2/ebitenutil" //Ajouter par moi pour l'importation d'images
)

// Draw affiche dans une image (en général, celle qui représente l'écran),
// la partie du sol qui est visible (qui doit avoir été calculée avec Get avant),
// UI string dans les paramètres de Draw correspond à g.PlayerOnUi.

// Variables de fonctionnement pour le rechargement des chunks lors du placement d'un bloc
var ReloadChunkX int
var ReloadChunkY int
var ReloadChunk bool

func (f Floor) Draw(screen *ebiten.Image, Particle *ebiten.Image, CharacterX int, CharacterY int, cameraY int, cameraX int, x_add int, y_add int, can_x bool, can_y bool, UI string) (PlaceBlocX, PlaceBlocY int, PlaceBloc bool, PlaceBlocID int) {
	// Cette fonction sert à dessiner les blocs sur la map

	// Ne pas faire attention aux tests juste en dessous, ils sont là lors d'une tentative de caméra fluide
	if can_x {
		x_add = 0
	}
	if can_y {
		y_add = 0
	}

	// Parcourir tous les blocs
	for ry := 0; ry < len(f.content); ry++ {
		for rx := 0; rx < len(f.content[0]); rx++ {

			x, y := rx, ry // Ne pas faire attention, on garde ceci pour la création de la caméra fluide
			mouseX, mouseY := ebiten.CursorPosition()
			var xGeneralPos int = mouseX / (configuration.Global.TileSize * 4)
			var yGeneralPos int = mouseY / (configuration.Global.TileSize * 4)
			configuration.Global.SelectPosX = cameraX/4 - (configuration.Global.RealNumTileX / 2) + xGeneralPos
			configuration.Global.SelectPosY = cameraY/4 - (configuration.Global.RealNumTileY / 2) + yGeneralPos

			// La condition si dessous vérifie si le bloc existe et qu'il ne s'agît pas de vide
			if f.bloc_at_pos(x, y) != -1 {

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(rx*configuration.Global.TileSize-x_add), float64(ry*configuration.Global.TileSize-y_add))

				//Gestion de l'éclairage
				op.ColorM.Scale(1.0, 1.0, 1.0, 1.0)
				lumi := f.LightAroundMe(x, y)
				op.ColorM.Translate(lumi, lumi, lumi, 0)

				shiftX := f.content[y][x] * configuration.Global.TileSize
				op.GeoM.Scale(4.0, 4.0)
				screen.DrawImage(assets.FloorImage.SubImage(
					image.Rect(shiftX, 0, shiftX+configuration.Global.TileSize, configuration.Global.TileSize),
				).(*ebiten.Image), op)

				//Affichage des particules
				if f.content[y][x] == configuration.Global.FurnaceID && configuration.Global.FireParticle && configuration.Global.ParticleOn {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Scale(4.0, 4.0)
					op.GeoM.Translate(float64(x*configuration.Global.TileSize*4+configuration.Global.TileSize+rand.Intn(configuration.Global.TileSize)), float64(y*configuration.Global.TileSize*4+configuration.Global.TileSize+rand.Intn(configuration.Global.TileSize)))
					screen.DrawImage(Particle, op)
				}

				// Rendre le bloc où la souris se positionne par dessus plus visible pour l'utilisateur
				if configuration.Global.CanSelect && UI == "" && x == xGeneralPos && y == yGeneralPos {
					selected, _, _ := ebitenutil.NewImageFromFile("../assets/Selected.png")
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(x*configuration.Global.TileSize), float64(y*configuration.Global.TileSize))
					op.GeoM.Scale(4.0, 4.0)
					screen.DrawImage(selected, op)

					// Poser des blocs (actuellement, seul le host peut poser des blocs et ils ne sont pas apparents chez les autres clients si le chunk à déjà été chargé)
					if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) && !configuration.Global.MapFileMode {
						NewX := cameraX + x - configuration.Global.RealNumTileX/2
						NewY := cameraY + y - configuration.Global.RealNumTileY/2
						fmt.Println(NewX, NewY)

						if configuration.Global.GameType == "join" { // Si il s'agit d'un client
							/* Désactivé car fonctionne mal du côté client
							*f.PlaceBlocX = NewX
							f.PlaceBlocY = NewY
							f.PlaceBlocID = 11
							f.PlaceBloc = true*/
						} else {
							// permet de poser un bloc (côté host)
							ReloadChunkX = DivisionPART(NewX)
							ReloadChunkY = DivisionPART(NewY)
							ReloadChunk = true
							SetBlocINActualChunkContent(0, NewX, NewY, configuration.Global.GameDir)
						}
					}
				}
			}
		}
	}
	return f.PlaceBlocX, f.PlaceBlocY, f.PlaceBloc, PlaceBlocID
}

func (f Floor) bloc_at_cam_pos(x, y int) int {
	// Permet d'obtenir le bloc qui se trouve à cam pos sans jamais sortir de la map (peut importe les coordonnées données en arguments

	if x >= 0 && y >= 0 && x < len(f.content[0]) && y < len(f.content) {
		return f.content[y][x]
	}
	return -1
}

func BlocisLight(blocID int) bool {
	// Détecte si le bloc <blocID> fais parmis des blocs considérés comme lumineux
	// Pour ajouter un bloc lumineux, il vous suffit de modifier dans le json configuration le paramètre << BlocsLights >>

	z := len(configuration.Global.BlocsLights)
	for i := 0; i < z; i++ {
		if blocID == configuration.Global.BlocsLights[i] {
			return true
		}
	}
	return false
}

func BlocisWall(blocID int) bool {
	// Détecte si le bloc <blocID> fais parmis des blocs considérés comme lumineux
	// Pour ajouter un bloc lumineux, il vous suffit de modifier dans le json configuration le paramètre << BlocsLights >>

	z := len(configuration.Global.Blocking)
	for i := 0; i < z; i++ {
		if blocID == configuration.Global.Blocking[i] {
			return true
		}
	}
	return false
}

func (f Floor) LightAroundMe(x, y int) float64 {
	// Cette fonction permet de trouver la lumineusité sur un bloc et à quel niveau

	// Elle vérifie si autour d'elle se trouve des blocs qui sont considérés comme des lumières en utilisant la fonction << BlocisLight() >>

	// Dans un premier temps, elle va vérifier elle même si le bloc à la position x y est lui même image.
	if BlocisLight(f.bloc_at_pos(x, y)) {
		return 0.0 // Ici puisque le bloc est effectivement une lumière, il ne perd rien en lumineusité
	}
	if configuration.Global.AmbiantLight < -0.3 {
		// Attention ! Il faut savoir que ça ne sert à rien de chercher de bloc émetteur de lumière si la lumière du << soleil >> est déjà
		// plus forte que ce qu'un bloc de lumière pourrait produire. D'où l'existe de la condition faite juste au dessus

		// Ici puisque le bloc n'est pas une lumière on va tout d'abord regarder à 1 case de distance de se blocs dans tous les sens possibles
		// si jamais il n'y aurais pas un bloc émetteur de lumière à côté de lui
		if BlocisLight(f.bloc_at_pos(x-1, y)) || BlocisLight(f.bloc_at_pos(x, y-1)) || BlocisLight(f.bloc_at_pos(x, y+1)) || BlocisLight(f.bloc_at_pos(x+1, y)) {
			return -0.3 // Ici oui donc on renvoie comme information que le bloc est éclairé comme lumière secondaire et donc à un niveau de lumière de -0.3
		}
		if configuration.Global.AmbiantLight < -0.4 {
			// Maintenant qu'on a vu que le bloc n'était pas éclairé à 1 bloc de distance, on va vérifier la même chose à 2 blocs de distance
			// On prend aussi en compte que la lumière ne peut pas passer derrière un bloc considéré comme << bloc mur >>
			//Par contre, si le bloc est lui même un mur, alors cette règle ne s'applique pas
			imAWall := BlocisWall(f.bloc_at_pos(x, y))
			if ((imAWall || !BlocisWall(f.bloc_at_pos(x-1, y))) && BlocisLight(f.bloc_at_pos(x-2, y))) ||
				((imAWall || !BlocisWall(f.bloc_at_pos(x, y-1))) && BlocisLight(f.bloc_at_pos(x, y-2))) ||
				((imAWall || !BlocisWall(f.bloc_at_pos(x, y+1))) && BlocisLight(f.bloc_at_pos(x, y+2))) ||
				((imAWall || !BlocisWall(f.bloc_at_pos(x+1, y))) && BlocisLight(f.bloc_at_pos(x+2, y))) ||
				((imAWall || !BlocisWall(f.bloc_at_pos(x-1, y))) || !BlocisWall(f.bloc_at_pos(x, y-1))) && BlocisLight(f.bloc_at_cam_pos(x-1, y-1)) ||
				((imAWall || !BlocisWall(f.bloc_at_pos(x+1, y)) || !BlocisWall(f.bloc_at_pos(x, y-1))) && BlocisLight(f.bloc_at_cam_pos(x+1, y-1))) ||
				((imAWall || !BlocisWall(f.bloc_at_pos(x-1, y)) || !BlocisWall(f.bloc_at_pos(x, y+1))) && BlocisLight(f.bloc_at_cam_pos(x-1, y+1))) ||
				((imAWall || !BlocisWall(f.bloc_at_pos(x+1, y)) || !BlocisWall(f.bloc_at_pos(x, y+1))) && BlocisLight(f.bloc_at_cam_pos(x+1, y+1))) {
				return -0.4 // Ici on a trouver un bloc lumière à 2 bloc de lui et donc on renvoie -4 en lumineusité
			}
		}
	}
	return configuration.Global.AmbiantLight // Ici puisque le bloc n'est pas éclairer, la dernière source de lumière possible c'est le soleil
}

func (f Floor) bloc_at_pos(x, y int) int {
	// Cette fonction sert aussi à obtenir le bloc de la cam à une certaine position mais est sujet à des modifications plus tard

	if x >= 0 && y >= 0 && x < len(f.content[0]) && y < len(f.content) {
		return f.content[y][x]
	}
	return -1
}

func DivisionPART(n int) int {
	if n < 0 && n%64 != 0 {
		return n/64 - 1
	}
	return n / 64
}

func DivisionPARTreste(n int) int {
	r := n % 64
	if r < 0 {
		r += 64
	}
	return r
}
