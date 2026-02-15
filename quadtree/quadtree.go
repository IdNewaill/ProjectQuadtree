package quadtree

import (
	"fmt"
	"strings"
)

// Quadtree est la structure de données pour les arbres
// quaternaires. Les champs non exportés sont :
//   - width, height : la taille en cases de la zone représentée
//     par l'arbre.
//   - root : le nœud qui est la racine de l'arbre.
type Quadtree struct {
	width, height int
	root          *node
}

// node représente un nœud d'arbre quaternaire. Les champs sont :
//   - topLeftX, topLeftY : les coordonnées (en cases) de la case
//     située en haut à gauche de la zone du terrain représentée
//     par ce nœud.
//   - width, height : la taille en cases de la zone représentée
//     par ce nœud.
//   - isLeaf : indique si le nœud est une feuille (true) ou pas (false).
//   - content : si le nœud est une feuille, indique le type de sol de
//     la zone représentée par ce nœud.
//   - xxxNode : Si le nœud n'est pas une feuille, donne une représentation
//     de la partie xxx de la zone représentée par ce nœud (dans ce cas, le
//     pointeur ne devrait jamais avoir pour valeur nil).
type node struct {
	topLeftX, topLeftY int
	width, height      int
	content            int
	isLeaf             bool
	topLeftNode        *node
	topRightNode       *node
	bottomLeftNode     *node
	bottomRightNode    *node
}

// GetAsText permet d'afficher un quadtree sous forme textuelle pour vérifier
// visuellement qu'il est correct (attention, pour des arbres un peu grands
// la visualisation devient peu lisible).
func (q Quadtree) GetAsText(name string) (asText string) {

	asText += fmt.Sprintf("Quadtree %s\n", name)
	asText += fmt.Sprintf("width: %d, height: %d\n", q.width, q.height)

	if q.root != nil {
		asText += q.root.getAsText(0, 7)
	} else {
		asText += "ERROR: quadtree has no root node"
	}

	return
}

// getAsText permet d'afficher un nœud sous forme textuelle, c'est une méthode
// utilitaire pour construire GetAsText.
func (n node) getAsText(shift, step int) (asText string) {

	shiftString := strings.Repeat(" ", shift)
	stepString := strings.Repeat(" ", step)

	asText += fmt.Sprintf("[node, top: (%d, %d), width: %d, height: %d,", n.topLeftX, n.topLeftY, n.width, n.height)

	if n.isLeaf {
		asText += fmt.Sprintf(" leaf with content %d]", n.content)
	} else {
		asText += fmt.Sprintf("\n%s%stopLeft:     %s", stepString, shiftString, n.topLeftNode.checkAndGetAsText(shift+step+13, step))
		asText += fmt.Sprintf("\n%s%stopRight:    %s", stepString, shiftString, n.topRightNode.checkAndGetAsText(shift+step+13, step))
		asText += fmt.Sprintf("\n%s%sbottomLeft:  %s", stepString, shiftString, n.bottomLeftNode.checkAndGetAsText(shift+step+13, step))
		asText += fmt.Sprintf("\n%s%sbottomRight: %s", stepString, shiftString, n.bottomRightNode.checkAndGetAsText(shift+step+13, step))
		asText += fmt.Sprintf("\n%s]", shiftString)
	}

	return
}

// checkAndGetAsText permet de vérifier qu'un nœud existe bien et
// d'afficher une erreur si ce n'est pas le cas. C'est une méthode
// utilitaire pour contruire getAsText.
func (n *node) checkAndGetAsText(shift, step int) (asText string) {

	if n == nil {
		asText = "ERROR: unexpected nil node"
		return
	}

	asText = n.getAsText(shift, step)
	return
}
