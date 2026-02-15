package quadtree

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}

func min(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}

func fillContentR(currentNode *node, contentHolderTopLeftX, contentHolderTopLeftY int, contentHolder [][]int) {

	contentHolderWidth := len(contentHolder[0])
	contentHolderHeight := len(contentHolder)

	interStartX := max(currentNode.topLeftX, contentHolderTopLeftX)
	interEndX := min(currentNode.topLeftX + currentNode.width, contentHolderTopLeftX + contentHolderWidth)
	interStartY := max(currentNode.topLeftY, contentHolderTopLeftY)
	interEndY := min(currentNode.topLeftY + currentNode.height, contentHolderTopLeftY + contentHolderHeight)

	// le nœud actuel n'interfaire pas avec contentHolder
	if interStartX >= interEndX || interStartY >= interEndY {
		return
	}

	if !currentNode.isLeaf {
		fillContentR(currentNode.topLeftNode, contentHolderTopLeftX, contentHolderTopLeftY, contentHolder)
		fillContentR(currentNode.topRightNode, contentHolderTopLeftX, contentHolderTopLeftY, contentHolder)
		fillContentR(currentNode.bottomLeftNode, contentHolderTopLeftX, contentHolderTopLeftY, contentHolder)
		fillContentR(currentNode.bottomRightNode, contentHolderTopLeftX, contentHolderTopLeftY, contentHolder)
		return
	}
	for x := interStartX; x < interEndX; x ++ {
		for y := interStartY; y < interEndY; y ++ {
			contentHolder[y -contentHolderTopLeftY][x - contentHolderTopLeftX] = currentNode.content
		}
	}
}

// GetContent remplit le tableau contentHolder (qui représente
// un terrain dont la case le plus en haut à gauche a pour coordonnées
// (topLeftX, topLeftY)) à partir du qadtree q.
func (q Quadtree) GetContent(topLeftX, topLeftY int, contentHolder [][]int) {
	fillContentR(q.root, topLeftX, topLeftY, contentHolder)
}
