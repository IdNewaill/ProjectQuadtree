package quadtree

// MakeFromArray construit un quadtree représentant un terrain
// étant donné un tableau représentant ce terrain.
func MakeFromArray(floorContent [][]int) (q Quadtree) {
	if len(floorContent) == 0 || len(floorContent[0]) == 0 {
		return Quadtree{} // Retourne un Quadtree vide si les données sont invalides
	}

	// Crée la racine de l'arbre à partir du tableau
	root := createNodeFromArray(floorContent, 0, 0, len(floorContent[0]), len(floorContent))
	return Quadtree{
		width:  len(floorContent[0]),
		height: len(floorContent),
		root:   root,
	}
}

// createNodeFromArray construit un nœud de l'arbre à partir d'une sous-zone du tableau.
func createNodeFromArray(content [][]int, x, y, width, height int) *node {
	if isHomogeneous(content, x, y, width, height) {
		// Si toutes les cases de la sous-zone ont la même valeur, crée une feuille
		return &node{
			topLeftX: x,
			topLeftY: y,
			width:    width,
			height:   height,
			isLeaf:   true,
			content:  content[y][x],
		}
	}

	// Divise la zone en 4 sous-zones et crée un nœud
	halfWidth := width / 2
	halfHeight := height / 2
	return &node{
		topLeftX:        x,
		topLeftY:        y,
		width:           width,
		height:          height,
		isLeaf:          false,
		topLeftNode:     createNodeFromArray(content, x, y, halfWidth, halfHeight),
		topRightNode:    createNodeFromArray(content, x+halfWidth, y, width-halfWidth, halfHeight),
		bottomLeftNode:  createNodeFromArray(content, x, y+halfHeight, halfWidth, height-halfHeight),
		bottomRightNode: createNodeFromArray(content, x+halfWidth, y+halfHeight, width-halfWidth, height-halfHeight),
	}
}

// isHomogeneous vérifie si toutes les cases dans une zone donnée ont la même valeur.
func isHomogeneous(content [][]int, x, y, width, height int) bool {
	refValue := content[y][x] // Valeur de référence (en haut à gauche de la zone)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if content[y+i][x+j] != refValue {
				return false // La zone n'est pas homogène
			}
		}
	}
	return true // La zone est homogène
}
