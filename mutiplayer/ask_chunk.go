package game_action

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/action_out"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
)

type AskChunk struct {
	X, Y     int
	Receiver string
}

func (a AskChunk) Execute(g *game.Game) {

	content := g.Floor.TryGetChunk(a.X, a.Y)
	if content != nil {
		fmt.Println("Using cached chunk", a.X, a.Y, "for", a.Receiver)
		g.ActionOutChan <- action_out.NewSetChunk(a.Receiver, a.X, a.Y, content)
		return
	}

	if configuration.Global.GameType != "join" {

		// charger le chunk en utilisant une nouvelle go routine pour ne pas ralentir le fonctionnement du jeu
		fmt.Println("Loading Chunk", a.X, a.Y, "for", a.Receiver)
		go game.LoadChunk(a.X, a.Y, a.Receiver, g.ActionOutChan, g.ChunkLoadedChan)
		return

	}

	// demander la mise à jour du chunk
	fmt.Println("Asking chunk", a.X, a.Y, "up")
	g.ActionOutChan <- action_out.NewAskChunk(
		"up",
		a.X, a.Y,
		a.Receiver+" "+configuration.Global.PlayerName,
	)

}

func decodeAskChunk(data []byte) (AskChunk, error) {

	action := AskChunk{}

	buff := *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&action)
	if err != nil {
		fmt.Println("connot decode ask chunk action:", err)
		return action, err
	}

	return action, nil

}
