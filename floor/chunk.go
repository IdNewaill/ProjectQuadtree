package floor

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"os"

	"gitlab.univ-nantes.fr/jezequel-l/quadtree/quadtree"

	"github.com/aquilax/go-perlin"
)

type Chunk struct {
	quadtree quadtree.Quadtree
}

func NewChunk(content [][]int) Chunk {
	return Chunk{
		quadtree: quadtree.MakeFromArray(content),
	}
}

func decodeContent(data []byte) ([][]int, error) {

	content := make([][]int, 64)
	for y := 0; y < 64; y++ {
		content[y] = make([]int, 64)
	}

	buff := *bytes.NewBuffer(data)
	dec := gob.NewDecoder(&buff)

	err := dec.Decode(&content)
	if err != nil {
		return content, errors.New("error decoding chunk content")
	}

	return content, nil
}

func encodeContent(content [][]int) ([]byte, error) {

	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(content)
	if err != nil {
		return buff.Bytes(), errors.New("error encoding chunk content")
	}

	return buff.Bytes(), nil
}

// Variables pour s'occuper du perlin noise
const noiseZoom = 20.0
const ilandSize = 100

var noise *perlin.Perlin

func genChunk(chunkX, chunkY int) [][]int {
	//Cette fonction permet de générer un chunck avec une taille de 64*64

	// Dans un premier temps, on vérifie que noise existe
	if noise == nil {
		panic("noise has not been init")
	}

	// Sinon on crée des tableaux de tableaux pour définir content
	content := make([][]int, 64)
	for y := 0; y < 64; y++ {
		content[y] = make([]int, 64)
	}

	// Puis ici on rempli content
	// Ici content est un chuck
	for x := 0; x < 64; x++ {
		for y := 0; y < 64; y++ {

			tileX := float64(chunkX*64 + x)
			tileY := float64(chunkY*64 + y)
			noiseValue := noise.Noise2D(tileX/noiseZoom, tileY/noiseZoom)

			distance := math.Sqrt(tileX*tileX + tileY*tileY)
			noiseValue -= (distance / ilandSize) * (distance / ilandSize)

			// Selon les valeur du perlin noise qui lui à été générer grâce à la seed donnée dans le fichier configuration.json
			// On choisi ici les blocs qui sont générés et qui correspondent à une tranche, par exemple entre 0 et 0.3
			// Et aussi avec un peu d'aléatoire
			if noiseValue < -0.8 {
				if rand.Intn(16) == 1 {
					content[y][x] = 24 // Bloc de glace un peu cassé
				} else if rand.Intn(10) == 1 {
					content[y][x] = 24 // Bloc de glace beaucoup cassé
				} else {
					content[y][x] = 23 // Bloc de glace
				}
			} else if noiseValue < -0.5 {
				if rand.Intn(30) == 1 {
					if rand.Intn(2) == 1 {
						content[y][x] = 12 + rand.Intn(2)
					} else {
						content[y][x] = 16 + rand.Intn(2)
					}
				} else {
					content[y][x] = 15 //bloc de sable
				}
			} else if noiseValue < 0.3 {
				if noiseValue < 0.3 && noiseValue > 0 {
					if rand.Intn(4) == 1 {
						content[y][x] = 18 + rand.Intn(4) // bloc d'herbe de fleurs
					} else if rand.Intn(4) == 1 {
						if rand.Intn(10) == 1 {
							content[y][x] = 14 // bloc de glowstone
						} else {
							content[y][x] = 4 // bloc de bois

							// Poser les feuilles des arbres
							trySetBloc(7, x-1, y, content) // A gauche
							trySetBloc(7, x, y-1, content) // En haut
							trySetBloc(7, x, y+1, content) // En bas
							trySetBloc(7, x+1, y, content) // A droite
						}
					} else {
						content[y][x] = 2 //bloc d'herbe
					}
				} else {
					if rand.Intn(16) == 1 {
						if rand.Intn(16) == 1 {
							content[y][x] = 8 + rand.Intn(3) //bloc de nourriture parmi citrouille, paille=>pain, pastèque
						} else {
							content[y][x] = 18 + rand.Intn(4) // bloc d'herbe de fleurs
						}
					} else {
						content[y][x] = 2 //bloc d'herbe
					}
				}
			} else {
				if rand.Intn(3) == 1 {
					if rand.Intn(2) == 1 {
						p := rand.Intn(3) // probabilité
						if p == 0 {
							content[y][x] = 3 // bloc de roche avec toile d'arraignée
						} else if p == 1 {
							content[y][x] = 31 // bloc de granite
						} else if p == 2 {
							content[y][x] = 28 // bloc d'andésite
						}

					} else if rand.Intn(2) == 1 {
						content[y][x] = 35 // minerais de charbon
					} else if rand.Intn(2) == 1 {
						content[y][x] = 37 // minerais de fer
					} else if rand.Intn(2) == 1 {
						content[y][x] = 38 // minerais d'or
					} else if rand.Intn(2) == 1 {
						content[y][x] = 39 // minerais de redstonne
					} else if rand.Intn(2) == 1 {
						content[y][x] = 36 // minerais de diamant
					}
				} else {
					content[y][x] = 32 // bloc de roche
				}
			}
		}
	}

	return content

}

func LoadChunkContent(x, y int, worldDir string) [][]int {

	var content [][]int

	filename := fmt.Sprintf("%s/%d-%d.chunk", worldDir, x, y)
	file, err := os.Open(filename)
	if err != nil {

		// si le fichier n'existe pas, alors c'est qu'on doit générer le chunk
		content = genChunk(x, y)

		// puis l'encoder
		data, err := encodeContent(content)
		if err != nil {
			panic(err)
		}

		// de façon à ensuite pouvoir le sauvegarder
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			fmt.Println(err)
			panic("Error writing chunk file")
		}

		// Puis renvoyer le chunk qui nous a été demandé
		return content

	}

	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		panic("Error reading chunk file")
	}

	// Read the file into a byte slice
	data := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(data)
	if err != nil {
		panic("Error reading chunk file")
	}

	// decode file content
	content, err = decodeContent(data)
	if err != nil {
		panic(err)
	}

	return content
}

func trySetBloc(blocID, x, y int, content [][]int) {
	// Cette fonction permet d'essayer de poser un bloc dans le chuck donné
	if x >= 0 && x < 64 && y >= 0 && y < 64 {
		content[y][x] = blocID
	}
}

func SetBlocINActualChunkContent(blocID, x, y int, worldDir string) {
	// Cette fonction permet de placer un bloc dans la map (Attention, elle ne permet pas d'actualiser la map)

	// Construire le nom du fichier de chunk
	filename := fmt.Sprintf("%s/%d-%d.chunk", worldDir, DivisionPART(x), DivisionPART(y))
	file, err := os.Open(filename) // Vérifier que le chunk existe
	if err == nil {
		defer file.Close()
		// Le chunk existe

		// Obtenir la taille du fichier
		stat, err := file.Stat()
		if err != nil {
			panic("Error reading chunk file")
		}

		// Transformation pour pouvoir comprendre le fichier
		data := make([]byte, stat.Size())
		_, err = bufio.NewReader(file).Read(data)
		if err != nil {
			panic("Error reading chunk file")
		}

		// Décoder le contenu du fichier
		var content [][]int
		content, err = decodeContent(data)
		if err != nil {
			panic(err)
		}

		// Placer le bloc si possible
		//trySetBloc(blocID, x%64, y%64, content)
		content[DivisionPARTreste(y)][DivisionPARTreste(x)] = blocID
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
