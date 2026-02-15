package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/assets"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/remote_player"
)

func main() {

	var configFileName string
	flag.StringVar(&configFileName, "config", "config.json", "select configuration file")

	flag.Parse()

	configuration.Load(configFileName)
	PlayerSpawnX := configuration.Global.PlayerSpawnX
	PlayerSpawnY := configuration.Global.PlayerSpawnY
	playerName := configuration.Global.PlayerName

	if strings.Contains(playerName, " ") {
		fmt.Println("Player name cannot contains spaces")
		return
	}

	// Ouvre des connexions à d’autres joueurs
	actionOutChan := make(chan game.ActionOut, 10)
	gameActionChan := make(chan game.GameAction, 10)
	newConnChan, connActionChan := remote_player.OpenConnections(gameActionChan, actionOutChan)

	// Se connecter si GameType == "join"
	if configuration.Global.GameType == "join" {

		addr := configuration.Global.ServerAddr
		isConnected := remote_player.Connect(
			addr,
			action_type.PlayerInitData{
				Name: playerName,
				X:    PlayerSpawnX,
				Y:    PlayerSpawnY,
			},
			gameActionChan,
			actionOutChan,
			newConnChan,
			connActionChan,
		)

		if !isConnected {
			fmt.Println("Connot connect to server.")
			return
		}

	} else if configuration.Global.GameType == "new" {

		// Créer game dir
		err := os.Mkdir(configuration.Global.GameDir, 0755)
		if err != nil {
			panic("Connot create game directory")
		}

	}

	// Démarrer le serveur de toute façon
	remote_player.StartServer(playerName, gameActionChan, actionOutChan, newConnChan, connActionChan)

	assets.Load()
	g := &game.Game{
		GameActionChan:  gameActionChan,
		ActionOutChan:   actionOutChan,
		ChunkLoadedChan: make(chan game.ChunkLoaded, 10),
	}
	g.Init(configuration.Global.Seed, PlayerSpawnX, PlayerSpawnY, playerName)

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	ebiten.SetWindowTitle("Quadtree")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
