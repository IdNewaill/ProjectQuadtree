package floor

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/aquilax/go-perlin"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/configuration"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/quadtree"
)

// Init initialise les structures de données internes de f.
func (f *Floor) Init(seed int64) {
	if configuration.Global.MapFileMode {
		f.content = make([][]int, configuration.Global.RealNumTileY)
		for y := 0; y < len(f.content); y++ {
			f.content[y] = make([]int, configuration.Global.RealNumTileX)
		}

		switch configuration.Global.FloorKind {
		case 1:
			f.NORMALfullContent = readFloorFromFile(configuration.Global.FloorFile)
		case 2:
			f.quadtreeContent = quadtree.MakeFromArray(readFloorFromFile(configuration.Global.FloorFile))
		}
	} else {
		// init perlin noise
		noise = perlin.NewPerlin(2, 2, 3, seed)

		// init fullContent map
		f.fullContent = make(map[Coords]*Chunk)

		// init content
		f.content = make([][]int, configuration.Global.NumTileY)
		for y := 0; y < len(f.content); y++ {
			f.content[y] = make([]int, configuration.Global.NumTileX)
		}
	}

}

// ################## SCRIPTS POUR MODE FICHIERS

// lecture du contenu d'un fichier représentant un terrain
// pour le stocker dans un tableau
func readFloorFromFile(fileName string) [][]int {
	// cette fonction permet de lire un fichier et de charger la map depuis ce fichier
	// elle prend en argument le nom du fichier

	readFile, err := os.Open(fileName)
	if err != nil { // vérifie si il n'y a pas eu de problème en ouvrant le fichier
		return nil
	}
	defer readFile.Close() // fermer le fichier dès la fin de la fonction

	fileScanner := bufio.NewScanner(readFile)
	var tableau [][]int

	for fileScanner.Scan() { // parcourt toutes les lignes de readFile
		line := fileScanner.Text()
		var br []int
		for _, char := range line {
			num, err := strconv.Atoi(string(char))
			if err != nil {
				fmt.Errorf("invalid character found in file: %v", err)
				panic("")
			}
			br = append(br, num)
		}
		tableau = append(tableau, br)
	}

	if err := fileScanner.Err(); err != nil {
		panic("Erreur lors de la lecture du fichier")
	}

	return tableau
}
