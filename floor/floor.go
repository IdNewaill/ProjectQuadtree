package floor

import (
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/quadtree"
)

type Coords struct {
	X, Y int
}

type Floor struct {
	content [][]int

	fullContent   map[Coords]*Chunk
	pendingChunks []Coords

	NORMALfullContent [][]int
	quadtreeContent   quadtree.Quadtree

	PlaceBlocX  int
	PlaceBlocY  int
	PlaceBlocID int
	PlaceBloc   bool
}

func (f *Floor) SetChunk(x, y int, chunk Chunk) {

	// supprimer le chunk s’il est en attente
	for k := range f.pendingChunks {
		if f.pendingChunks[k].X == x && f.pendingChunks[k].Y == y {
			f.pendingChunks[k] = f.pendingChunks[len(f.pendingChunks)-1]
			f.pendingChunks = f.pendingChunks[:len(f.pendingChunks)-1]
			break
		}
	}

	f.fullContent[Coords{x, y}] = &chunk
}

func (f *Floor) TryGetChunk(x, y int) [][]int {

	// vérifie si le chunk n'est pas déjà stocker quelque part en mémoire
	chunk, ok := f.fullContent[Coords{x, y}]
	if !ok {
		return nil
	}

	// crée un content vide
	content := make([][]int, 64)
	for y := 0; y < 64; y++ {
		content[y] = make([]int, 64)
	}

	chunk.quadtree.GetContent(0, 0, content)
	f.content = content
	return content
}

func (f Floor) GetContent() [][]int { // Fonction Ajoutée par moi
	return f.content
}
