package game_action

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/floor"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/game"
)

type SetChunk struct {
	X, Y    int
	Content [][]int
}

func (a SetChunk) Execute(g *game.Game) {
	// Cette fonction agit quand un client demande de remplacer un des chunk actuels par un nouveau

	fmt.Println("Received chunk :", a.X, a.Y)
	if GetDifferenceBetweenContent(g.Floor.TryGetChunk(a.X, a.Y), a.Content) {
		if configuration.Global.GameType != "join" {
			SetActualChunkContent(a.Content, a.X, a.Y, configuration.Global.GameDir)
			//g.Floor.SetChunk(a.X, a.Y, floor.NewChunk(a.Content))
		}
		g.Floor.SetChunk(a.X, a.Y, floor.NewChunk(a.Content))
	}

}

func decodeSetChunk(data []byte) (SetChunk, error) {
	// Cette fonction permet de décoder un chunk

	action := SetChunk{}

	buff := *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&action)
	if err != nil {
		fmt.Println("connot decode set chunk action:", err)
		return action, err
	}

	return action, nil

}

func SetActualChunkContent(content [][]int, x, y int, worldDir string) {
	// Cette fonction permet de placer un bloc dans la map (Attention, elle ne permet pas d'actualiser la map)

	// D'abord on va vérifier que le fichier chunk à la position x y existe
	// Pour cela, il faut d'abord construire le nom du fichier de chunk
	// son nom est en fait créer à partir des coordonnées
	filename := fmt.Sprintf("%s/%d-%d.chunk", worldDir, x, y)
	file, err := os.Open(filename) // Vérifier que le chunk existe
	if err == nil {
		defer file.Close()
		// Réencoder puis fermer le mode lecture
		encodedContent, err := encodeContent(content)
		if err != nil {
			panic(err)
		}

		// Ouvrir le fichier en mode écriture
		saveFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic("Error saving chunk file")
		}
		defer saveFile.Close()

		// Sauvegarder le changement
		_, err = saveFile.Write(encodedContent)
		if err != nil {
			panic("Error writing to chunk file")
		}
	} else {
		// Si le fichier n'existe pas, déclencher une erreur
		panic("Error, chunk does not exist")
	}
}

func encodeContent(content [][]int) ([]byte, error) {
	// Cette fonction permet d'encoder content pour le quadtree avec perlin noise (notamment pour le multijoueur)

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(content)
	if err != nil {
		return buff.Bytes(), errors.New("error encoding chunk content")
	}

	return buff.Bytes(), nil
}

func GetDifferenceBetweenContent(contentA, contentB [][]int) bool {
	// Cette fonction était utilisée pour une autre fonction qui elle permettait d'update le terrain après qu'un bloc ai
	// été posé. Cette fonction sert essentiellement à vérifier qu'un joueur n'essaye pas de détruire tout un chunk.

	// Elle vérifie si il y a moins de 10 différences entre le content actuel et reçu
	err := 0
	for y := 0; y < len(contentA); y++ {
		for x := 0; x < len(contentA); x++ {
			if contentA[y][x] != contentB[y][x] {
				err += 1
				if err > 9 {
					return false
				}
			}
		}
	}
	return true
}
