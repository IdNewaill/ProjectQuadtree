package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_out"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
)

// Update met à jour les données du jeu à chaque 1/60 de seconde.
// Il faut bien faire attention à l'ordre des mises-à-jour car elles
// dépendent les unes des autres (par exemple, pour le moment, la
// mise-à-jour de la caméra dépend de celle du personnage et la définition
// du terrain dépend de celle de la caméra).

// Variable de test
//var bloc_pose bool

func (g *Game) Update() error {
	// Permet d'ouvrir ou fermer l'inventaire
	if inpututil.IsKeyJustPressed(ebiten.KeyE) && !configuration.Global.OriginalKey {
		configuration.Global.DebugMode = false
		configuration.Global.UiInventaireOn = !configuration.Global.UiInventaireOn
		if configuration.Global.UiInventaireOn {
			g.PlayerOnUi = "inventory"
		} else {
			g.PlayerOnUi = ""
		}
	}

	// Permettre au joueur d'appuyer sur certaines touches seulement si il ne se trouve pas dans une interface
	if g.PlayerOnUi == "" && !configuration.Global.OriginalKey {
		if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
			configuration.Global.DebugMode = !configuration.Global.DebugMode
		}

		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			configuration.Global.HotbarSelect = 0
		}
		if inpututil.IsKeyJustPressed(ebiten.Key2) {
			configuration.Global.HotbarSelect = 1
		}
		if inpututil.IsKeyJustPressed(ebiten.Key3) {
			configuration.Global.HotbarSelect = 2
		}
		if inpututil.IsKeyJustPressed(ebiten.Key4) {
			configuration.Global.HotbarSelect = 3
		}
		if inpututil.IsKeyJustPressed(ebiten.Key5) {
			configuration.Global.HotbarSelect = 4
		}
		if inpututil.IsKeyJustPressed(ebiten.Key6) {
			configuration.Global.HotbarSelect = 5
		}
		if inpututil.IsKeyJustPressed(ebiten.Key7) {
			configuration.Global.HotbarSelect = 6
		}
		if inpututil.IsKeyJustPressed(ebiten.Key8) {
			configuration.Global.HotbarSelect = 7
		}
		if inpututil.IsKeyJustPressed(ebiten.Key9) {
			configuration.Global.HotbarSelect = 8
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) && configuration.Global.OriginalKey {
		configuration.Global.DebugMode = !configuration.Global.DebugMode
	}

	blocking := g.Floor.Blocking(g.Character.X, g.Character.Y, g.Camera.X, g.Camera.Y)

	prevOrientation := g.Character.Orientation
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if g.Character.TryMoveUp(blocking) {
			g.ActionOutChan <- action_out.NewNotifyPlayerMove(
				"all",
				configuration.Global.PlayerName,
				g.Character.X,
				g.Character.Y-1,
			)
		} else {
			if prevOrientation != g.Character.Orientation {
				g.ActionOutChan <- action_out.NewNotifyOrientationChange(
					"all",
					configuration.Global.PlayerName,
					g.Character.Orientation,
				)
			}
		}

	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		fmt.Println("Passer par game update.go 98")
		if g.Character.TryMoveRight(blocking) {
			fmt.Println("AAA")
			g.ActionOutChan <- action_out.NewNotifyPlayerMove(
				"all",
				configuration.Global.PlayerName,
				g.Character.X+1,
				g.Character.Y,
			)
		} else {
			if prevOrientation != g.Character.Orientation {
				g.ActionOutChan <- action_out.NewNotifyOrientationChange(
					"all",
					configuration.Global.PlayerName,
					g.Character.Orientation,
				)
			}
		}

	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {

		if g.Character.TryMoveDown(blocking) {

			g.ActionOutChan <- action_out.NewNotifyPlayerMove(
				"all",
				configuration.Global.PlayerName,
				g.Character.X,
				g.Character.Y+1,
			)
		} else {
			if prevOrientation != g.Character.Orientation {
				g.ActionOutChan <- action_out.NewNotifyOrientationChange(
					"all",
					configuration.Global.PlayerName,
					g.Character.Orientation,
				)
			}
		}

	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {

		if g.Character.TryMoveLeft(blocking) {
			g.ActionOutChan <- action_out.NewNotifyPlayerMove(
				"all",
				configuration.Global.PlayerName,
				g.Character.X-1,
				g.Character.Y,
			)
		} else {
			if prevOrientation != g.Character.Orientation {
				g.ActionOutChan <- action_out.NewNotifyOrientationChange(
					"all",
					configuration.Global.PlayerName,
					g.Character.Orientation,
				)
			}
		}

	}

	// S'occuper du blocage de joueur et autre
	x_add, y_add := g.Character.GetContentPLayerAddToCamPos()
	moving := g.Character.Update(blocking) // moving := g.Character.Update(blocking)
	for index := range g.OtherPlayers {
		g.OtherPlayers[index].Update([6]bool{true, true, true, true, true, true})
	}

	g.Camera.Update(g.Character.X, g.Character.Y, x_add, y_add)

	// S'occuper des sons

	if configuration.Global.SoundsON {
		g.manage_sounds(moving)
	}

	// charger un chunk
	chunkToLoad := g.Floor.Update(g.Camera.X, g.Camera.Y)

	if !configuration.Global.MapFileMode {
		// si le joueur est en mode << join >> autrement dis qu'il est dans la partie de quelqu'un d'autre
		if configuration.Global.GameType == "join" {
			// Pour permettre à un client de poser un bloc
			if g.PlaceBloc && configuration.Global.CanPlace {
				Selected := configuration.Global.HotbarSelect
				blocID := configuration.Global.InventoryItems[Selected][0]
				if configuration.Global.InventoryItems[Selected][1] > 0 && blocID != -1 {
					bloc_pose := g.clientPoseUnBloc("all", g.PlaceBlocX, g.PlaceBlocY, blocID)
					g.PlaceBloc = !bloc_pose
				}
			}

			for i := range chunkToLoad {
				fmt.Println("asking chunk :", chunkToLoad[i].X, chunkToLoad[i].Y)
				g.ActionOutChan <- action_out.NewAskChunk(
					"up",
					chunkToLoad[i].X,
					chunkToLoad[i].Y,
					configuration.Global.PlayerName,
				)
			}

		} else {
			// charger les chunks
			for i := range chunkToLoad {
				go LoadChunk(chunkToLoad[i].X, chunkToLoad[i].Y, "", g.ActionOutChan, g.ChunkLoadedChan)
			}
		}
		// executer les actions reçues
		for {
			select {
			case action := <-g.GameActionChan:
				fmt.Println("hi")
				action.Execute(g)
			case chunk := <-g.ChunkLoadedChan:
				g.Floor.SetChunk(chunk.coord.X, chunk.coord.Y, chunk.chunk)
			default:
				return nil
			}
		}
	}
	return nil
}

func (g *Game) clientPoseUnBloc(sendTo string, blocPoseX, blocPoseY int, blocID int) bool {
	var content [][]int
	ChunkX, ChunkY := DivisionPART(blocPoseX), DivisionPART(blocPoseY)
	MouseX, MouseY := DivisionPARTreste(blocPoseX), DivisionPARTreste(blocPoseY)
	content = g.Floor.TryGetChunk(ChunkX, ChunkY)
	if len(content) > 0 {
		if content[MouseY][MouseX] != blocID {
			content[MouseY][MouseX] = blocID
			g.ActionOutChan <- action_out.NewSetChunk(sendTo, ChunkX, ChunkY, content)
		}
		return true
	}
	return false
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
