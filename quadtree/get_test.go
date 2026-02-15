package quadtree

import "testing"

type requestedWindow struct {
	topLeftX, topLeftY int
	width, height      int
	expectedResult     [][]int
}

// Vérife la fonction GetContent avec des quadtrees où les
// réductions dans le cas où toute la partie du terrain
// représentée par un nœud n'a qu'un seul type de sol n'ont
// pas été effectuées.
func TestGetContentNoReduction(t *testing.T) {

	trees := []Quadtree{
		{width: 2, height: 2,
			root: &node{topLeftX: 0, topLeftY: 0, width: 2, height: 2,
				topLeftNode:     &node{topLeftX: 0, topLeftY: 0, width: 1, height: 1, isLeaf: true, content: 1},
				topRightNode:    &node{topLeftX: 1, topLeftY: 0, width: 1, height: 1, isLeaf: true, content: 2},
				bottomLeftNode:  &node{topLeftX: 0, topLeftY: 1, width: 1, height: 1, isLeaf: true, content: 3},
				bottomRightNode: &node{topLeftX: 1, topLeftY: 1, width: 1, height: 1, isLeaf: true, content: 4},
			},
		},
		{width: 6, height: 5,
			root: &node{topLeftX: 0, topLeftY: 0, width: 6, height: 5,
				topLeftNode: &node{topLeftX: 0, topLeftY: 0, width: 3, height: 2,
					topLeftNode: &node{topLeftX: 0, topLeftY: 0, width: 1, height: 1, isLeaf: true, content: 1},
					topRightNode: &node{topLeftX: 1, topLeftY: 0, width: 2, height: 1,
						topLeftNode:     &node{topLeftX: 1, topLeftY: 0, width: 1, height: 0, isLeaf: true, content: 0},
						topRightNode:    &node{topLeftX: 2, topLeftY: 0, width: 1, height: 0, isLeaf: true, content: 0},
						bottomLeftNode:  &node{topLeftX: 1, topLeftY: 0, width: 1, height: 1, isLeaf: true, content: 2},
						bottomRightNode: &node{topLeftX: 2, topLeftY: 0, width: 1, height: 1, isLeaf: true, content: 3},
					},
					bottomLeftNode: &node{topLeftX: 0, topLeftY: 1, width: 1, height: 1, isLeaf: true, content: 2},
					bottomRightNode: &node{topLeftX: 1, topLeftY: 1, width: 2, height: 1,
						topLeftNode:     &node{topLeftX: 1, topLeftY: 1, width: 1, height: 0, isLeaf: true, content: 0},
						topRightNode:    &node{topLeftX: 2, topLeftY: 1, width: 1, height: 0, isLeaf: true, content: 0},
						bottomLeftNode:  &node{topLeftX: 1, topLeftY: 1, width: 1, height: 1, isLeaf: true, content: 3},
						bottomRightNode: &node{topLeftX: 2, topLeftY: 1, width: 1, height: 1, isLeaf: true, content: 4},
					},
				},
				topRightNode: &node{topLeftX: 3, topLeftY: 0, width: 3, height: 2,
					topLeftNode: &node{topLeftX: 3, topLeftY: 0, width: 1, height: 1, isLeaf: true, content: 4},
					topRightNode: &node{topLeftX: 4, topLeftY: 0, width: 2, height: 1,
						topLeftNode:     &node{topLeftX: 4, topLeftY: 0, width: 1, height: 0, isLeaf: true, content: 0},
						topRightNode:    &node{topLeftX: 5, topLeftY: 0, width: 1, height: 0, isLeaf: true, content: 0},
						bottomLeftNode:  &node{topLeftX: 4, topLeftY: 0, width: 1, height: 1, isLeaf: true, content: 1},
						bottomRightNode: &node{topLeftX: 5, topLeftY: 0, width: 1, height: 1, isLeaf: true, content: 2},
					},
					bottomLeftNode: &node{topLeftX: 3, topLeftY: 1, width: 1, height: 1, isLeaf: true, content: 1},
					bottomRightNode: &node{topLeftX: 4, topLeftY: 1, width: 2, height: 1,
						topLeftNode:     &node{topLeftX: 4, topLeftY: 1, width: 1, height: 0, isLeaf: true, content: 0},
						topRightNode:    &node{topLeftX: 5, topLeftY: 1, width: 1, height: 0, isLeaf: true, content: 0},
						bottomLeftNode:  &node{topLeftX: 4, topLeftY: 1, width: 1, height: 1, isLeaf: true, content: 2},
						bottomRightNode: &node{topLeftX: 5, topLeftY: 1, width: 1, height: 1, isLeaf: true, content: 3},
					},
				},
				bottomLeftNode: &node{topLeftX: 0, topLeftY: 2, width: 3, height: 3,
					topLeftNode: &node{topLeftX: 0, topLeftY: 2, width: 1, height: 1, isLeaf: true, content: 3},
					topRightNode: &node{topLeftX: 1, topLeftY: 2, width: 2, height: 1,
						topLeftNode:     &node{topLeftX: 1, topLeftY: 2, width: 1, height: 0, isLeaf: true, content: 0},
						topRightNode:    &node{topLeftX: 2, topLeftY: 2, width: 1, height: 0, isLeaf: true, content: 0},
						bottomLeftNode:  &node{topLeftX: 1, topLeftY: 2, width: 1, height: 1, isLeaf: true, content: 4},
						bottomRightNode: &node{topLeftX: 2, topLeftY: 2, width: 1, height: 1, isLeaf: true, content: 1},
					},
					bottomLeftNode: &node{topLeftX: 0, topLeftY: 3, width: 1, height: 2,
						topLeftNode:     &node{topLeftX: 0, topLeftY: 3, width: 0, height: 1, isLeaf: true, content: 0},
						topRightNode:    &node{topLeftX: 0, topLeftY: 3, width: 1, height: 1, isLeaf: true, content: 4},
						bottomLeftNode:  &node{topLeftX: 0, topLeftY: 4, width: 0, height: 1, isLeaf: true, content: 0},
						bottomRightNode: &node{topLeftX: 0, topLeftY: 4, width: 1, height: 1, isLeaf: true, content: 1},
					},
					bottomRightNode: &node{topLeftX: 1, topLeftY: 3, width: 2, height: 2,
						topLeftNode:     &node{topLeftX: 1, topLeftY: 3, width: 1, height: 1, isLeaf: true, content: 1},
						topRightNode:    &node{topLeftX: 2, topLeftY: 3, width: 1, height: 1, isLeaf: true, content: 2},
						bottomLeftNode:  &node{topLeftX: 1, topLeftY: 4, width: 1, height: 1, isLeaf: true, content: 2},
						bottomRightNode: &node{topLeftX: 2, topLeftY: 4, width: 1, height: 1, isLeaf: true, content: 3},
					},
				},
				bottomRightNode: &node{topLeftX: 3, topLeftY: 2, width: 3, height: 3,
					topLeftNode: &node{topLeftX: 3, topLeftY: 2, width: 1, height: 1, isLeaf: true, content: 2},
					topRightNode: &node{topLeftX: 4, topLeftY: 2, width: 2, height: 1,
						topLeftNode:     &node{topLeftX: 4, topLeftY: 2, width: 1, height: 0, isLeaf: true, content: 0},
						topRightNode:    &node{topLeftX: 5, topLeftY: 2, width: 1, height: 0, isLeaf: true, content: 0},
						bottomLeftNode:  &node{topLeftX: 4, topLeftY: 2, width: 1, height: 1, isLeaf: true, content: 3},
						bottomRightNode: &node{topLeftX: 5, topLeftY: 2, width: 1, height: 1, isLeaf: true, content: 4},
					},
					bottomLeftNode: &node{topLeftX: 3, topLeftY: 3, width: 1, height: 2,
						topLeftNode:     &node{topLeftX: 3, topLeftY: 3, width: 0, height: 1, isLeaf: true, content: 0},
						topRightNode:    &node{topLeftX: 3, topLeftY: 3, width: 1, height: 1, isLeaf: true, content: 3},
						bottomLeftNode:  &node{topLeftX: 3, topLeftY: 4, width: 0, height: 1, isLeaf: true, content: 0},
						bottomRightNode: &node{topLeftX: 3, topLeftY: 4, width: 1, height: 1, isLeaf: true, content: 4},
					},
					bottomRightNode: &node{topLeftX: 4, topLeftY: 3, width: 2, height: 2,
						topLeftNode:     &node{topLeftX: 4, topLeftY: 3, width: 1, height: 1, isLeaf: true, content: 4},
						topRightNode:    &node{topLeftX: 5, topLeftY: 3, width: 1, height: 1, isLeaf: true, content: 1},
						bottomLeftNode:  &node{topLeftX: 4, topLeftY: 4, width: 1, height: 1, isLeaf: true, content: 1},
						bottomRightNode: &node{topLeftX: 5, topLeftY: 4, width: 1, height: 1, isLeaf: true, content: 2},
					},
				},
			},
		},
	}

	requests := [][]requestedWindow{
		{
			{topLeftX: 0, topLeftY: 0, width: 2, height: 2, expectedResult: [][]int{{1, 2}, {3, 4}}},
			{topLeftX: 0, topLeftY: 0, width: 1, height: 1, expectedResult: [][]int{{1}}},
			{topLeftX: 1, topLeftY: 0, width: 1, height: 1, expectedResult: [][]int{{2}}},
			{topLeftX: 0, topLeftY: 1, width: 1, height: 1, expectedResult: [][]int{{3}}},
			{topLeftX: 1, topLeftY: 1, width: 1, height: 1, expectedResult: [][]int{{4}}},
			{topLeftX: 0, topLeftY: 0, width: 2, height: 1, expectedResult: [][]int{{1, 2}}},
			{topLeftX: 0, topLeftY: 1, width: 2, height: 1, expectedResult: [][]int{{3, 4}}},
			{topLeftX: 0, topLeftY: 0, width: 1, height: 2, expectedResult: [][]int{{1}, {3}}},
			{topLeftX: 1, topLeftY: 0, width: 1, height: 2, expectedResult: [][]int{{2}, {4}}},
			{topLeftX: 1, topLeftY: 0, width: 2, height: 2, expectedResult: [][]int{{2, -1}, {4, -1}}},
			{topLeftX: 0, topLeftY: 1, width: 2, height: 2, expectedResult: [][]int{{3, 4}, {-1, -1}}},
			{topLeftX: 1, topLeftY: 1, width: 2, height: 2, expectedResult: [][]int{{4, -1}, {-1, -1}}},
			{topLeftX: 2, topLeftY: 2, width: 2, height: 2, expectedResult: [][]int{{-1, -1}, {-1, -1}}},
			{topLeftX: -1, topLeftY: 0, width: 2, height: 2, expectedResult: [][]int{{-1, 1}, {-1, 3}}},
			{topLeftX: 0, topLeftY: -1, width: 2, height: 2, expectedResult: [][]int{{-1, -1}, {1, 2}}},
			{topLeftX: -1, topLeftY: -1, width: 2, height: 2, expectedResult: [][]int{{-1, -1}, {-1, 1}}},
			{topLeftX: -2, topLeftY: -2, width: 2, height: 2, expectedResult: [][]int{{-1, -1}, {-1, -1}}},
		},
		{
			{topLeftX: 0, topLeftY: 0, width: 6, height: 5, expectedResult: [][]int{
				{1, 2, 3, 4, 1, 2},
				{2, 3, 4, 1, 2, 3},
				{3, 4, 1, 2, 3, 4},
				{4, 1, 2, 3, 4, 1},
				{1, 2, 3, 4, 1, 2},
			}},
			{topLeftX: 0, topLeftY: 0, width: 4, height: 4, expectedResult: [][]int{
				{1, 2, 3, 4},
				{2, 3, 4, 1},
				{3, 4, 1, 2},
				{4, 1, 2, 3},
			}},
			{topLeftX: 1, topLeftY: 0, width: 4, height: 4, expectedResult: [][]int{
				{2, 3, 4, 1},
				{3, 4, 1, 2},
				{4, 1, 2, 3},
				{1, 2, 3, 4},
			}},
			{topLeftX: 2, topLeftY: 0, width: 4, height: 4, expectedResult: [][]int{
				{3, 4, 1, 2},
				{4, 1, 2, 3},
				{1, 2, 3, 4},
				{2, 3, 4, 1},
			}},
			{topLeftX: 0, topLeftY: 1, width: 4, height: 4, expectedResult: [][]int{
				{2, 3, 4, 1},
				{3, 4, 1, 2},
				{4, 1, 2, 3},
				{1, 2, 3, 4},
			}},
			{topLeftX: 1, topLeftY: 1, width: 4, height: 4, expectedResult: [][]int{
				{3, 4, 1, 2},
				{4, 1, 2, 3},
				{1, 2, 3, 4},
				{2, 3, 4, 1},
			}},
			{topLeftX: 2, topLeftY: 1, width: 4, height: 4, expectedResult: [][]int{
				{4, 1, 2, 3},
				{1, 2, 3, 4},
				{2, 3, 4, 1},
				{3, 4, 1, 2},
			}},
			{topLeftX: 4, topLeftY: 2, width: 4, height: 4, expectedResult: [][]int{
				{3, 4, -1, -1},
				{4, 1, -1, -1},
				{1, 2, -1, -1},
				{-1, -1, -1, -1},
			}},
			{topLeftX: -2, topLeftY: -1, width: 4, height: 4, expectedResult: [][]int{
				{-1, -1, -1, -1},
				{-1, -1, 1, 2},
				{-1, -1, 2, 3},
				{-1, -1, 3, 4},
			}},
			{topLeftX: -1, topLeftY: 3, width: 8, height: 3, expectedResult: [][]int{
				{-1, 4, 1, 2, 3, 4, 1, -1},
				{-1, 1, 2, 3, 4, 1, 2, -1},
				{-1, -1, -1, -1, -1, -1, -1, -1},
			}},
		},
	}

	for i, rSet := range requests {

		for _, r := range rSet {

			res := make([][]int, r.height)
			for y := range res {
				res[y] = make([]int, r.width)
			}

			trees[i].GetContent(r.topLeftX, r.topLeftY, res)

			if !equalContent(r.expectedResult, res) {
				treeAsText := trees[i].GetAsText("fullContent")
				t.Fatalf("Le résultat attendu pour le quadtree \n%s\n avec (topLeftX, topLeftY) = (%d, %d) était le tableau %v mais le résultat retourné par GetContent est %v", treeAsText, r.topLeftX, r.topLeftY, r.expectedResult, res)
			}
		}

	}

}

// Vérife la fonction GetContent avec des quadtrees où les
// réductions dans le cas où toute la partie du terrain
// représentée par un nœud n'a qu'un seul type de sol ont
// été effectuées.
func TestGetContentReduction(t *testing.T) {

	trees := []Quadtree{
		{width: 8, height: 8,
			root: &node{topLeftX: 0, topLeftY: 0, width: 8, height: 8,
				topLeftNode: &node{topLeftX: 0, topLeftY: 0, width: 4, height: 4,
					topLeftNode:     &node{topLeftX: 0, topLeftY: 0, width: 2, height: 2, isLeaf: true, content: 1},
					topRightNode:    &node{topLeftX: 2, topLeftY: 0, width: 2, height: 2, isLeaf: true, content: 2},
					bottomLeftNode:  &node{topLeftX: 0, topLeftY: 2, width: 2, height: 2, isLeaf: true, content: 2},
					bottomRightNode: &node{topLeftX: 2, topLeftY: 2, width: 2, height: 2, isLeaf: true, content: 3},
				},
				topRightNode: &node{topLeftX: 4, topLeftY: 0, width: 4, height: 4,
					topLeftNode:     &node{topLeftX: 4, topLeftY: 0, width: 2, height: 2, isLeaf: true, content: 3},
					topRightNode:    &node{topLeftX: 6, topLeftY: 0, width: 2, height: 2, isLeaf: true, content: 4},
					bottomLeftNode:  &node{topLeftX: 4, topLeftY: 2, width: 2, height: 2, isLeaf: true, content: 4},
					bottomRightNode: &node{topLeftX: 6, topLeftY: 2, width: 2, height: 2, isLeaf: true, content: 1},
				},
				bottomLeftNode:  &node{topLeftX: 0, topLeftY: 4, width: 4, height: 4, isLeaf: true, content: 1},
				bottomRightNode: &node{topLeftX: 4, topLeftY: 4, width: 4, height: 4, isLeaf: true, content: 2},
			},
		},
		{width: 6, height: 5,
			root: &node{topLeftX: 0, topLeftY: 0, width: 6, height: 5,
				topLeftNode:  &node{topLeftX: 0, topLeftY: 0, width: 3, height: 2, isLeaf: true, content: 1},
				topRightNode: &node{topLeftX: 3, topLeftY: 0, width: 3, height: 2, isLeaf: true, content: 2},
				bottomLeftNode: &node{topLeftX: 0, topLeftY: 2, width: 3, height: 3,
					topLeftNode:     &node{topLeftX: 0, topLeftY: 2, width: 1, height: 1, isLeaf: true, content: 1},
					topRightNode:    &node{topLeftX: 1, topLeftY: 2, width: 2, height: 1, isLeaf: true, content: 1},
					bottomLeftNode:  &node{topLeftX: 0, topLeftY: 3, width: 1, height: 2, isLeaf: true, content: 3},
					bottomRightNode: &node{topLeftX: 1, topLeftY: 3, width: 2, height: 2, isLeaf: true, content: 3},
				},
				bottomRightNode: &node{topLeftX: 3, topLeftY: 2, width: 3, height: 3, isLeaf: true, content: 4},
			},
		},
	}

	requests := [][]requestedWindow{
		{
			{topLeftX: 0, topLeftY: 0, width: 8, height: 8, expectedResult: [][]int{
				{1, 1, 2, 2, 3, 3, 4, 4},
				{1, 1, 2, 2, 3, 3, 4, 4},
				{2, 2, 3, 3, 4, 4, 1, 1},
				{2, 2, 3, 3, 4, 4, 1, 1},
				{1, 1, 1, 1, 2, 2, 2, 2},
				{1, 1, 1, 1, 2, 2, 2, 2},
				{1, 1, 1, 1, 2, 2, 2, 2},
				{1, 1, 1, 1, 2, 2, 2, 2},
			}},
			{topLeftX: 0, topLeftY: 0, width: 5, height: 5, expectedResult: [][]int{
				{1, 1, 2, 2, 3},
				{1, 1, 2, 2, 3},
				{2, 2, 3, 3, 4},
				{2, 2, 3, 3, 4},
				{1, 1, 1, 1, 2},
			}}, {topLeftX: 3, topLeftY: 3, width: 5, height: 5, expectedResult: [][]int{
				{3, 4, 4, 1, 1},
				{1, 2, 2, 2, 2},
				{1, 2, 2, 2, 2},
				{1, 2, 2, 2, 2},
				{1, 2, 2, 2, 2},
			}},
			{topLeftX: 1, topLeftY: 1, width: 6, height: 6, expectedResult: [][]int{
				{1, 2, 2, 3, 3, 4},
				{2, 3, 3, 4, 4, 1},
				{2, 3, 3, 4, 4, 1},
				{1, 1, 1, 2, 2, 2},
				{1, 1, 1, 2, 2, 2},
				{1, 1, 1, 2, 2, 2},
			}},
			{topLeftX: 5, topLeftY: 2, width: 7, height: 7, expectedResult: [][]int{
				{4, 1, 1, -1, -1, -1, -1},
				{4, 1, 1, -1, -1, -1, -1},
				{2, 2, 2, -1, -1, -1, -1},
				{2, 2, 2, -1, -1, -1, -1},
				{2, 2, 2, -1, -1, -1, -1},
				{2, 2, 2, -1, -1, -1, -1},
				{-1, -1, -1, -1, -1, -1, -1},
			}},
		},
		{
			{topLeftX: 0, topLeftY: 0, width: 6, height: 5, expectedResult: [][]int{
				{1, 1, 1, 2, 2, 2},
				{1, 1, 1, 2, 2, 2},
				{1, 1, 1, 4, 4, 4},
				{3, 3, 3, 4, 4, 4},
				{3, 3, 3, 4, 4, 4},
			}},
			{topLeftX: 0, topLeftY: 0, width: 4, height: 4, expectedResult: [][]int{
				{1, 1, 1, 2},
				{1, 1, 1, 2},
				{1, 1, 1, 4},
				{3, 3, 3, 4},
			}},
			{topLeftX: 1, topLeftY: 0, width: 4, height: 4, expectedResult: [][]int{
				{1, 1, 2, 2},
				{1, 1, 2, 2},
				{1, 1, 4, 4},
				{3, 3, 4, 4},
			}},
			{topLeftX: 2, topLeftY: 0, width: 4, height: 4, expectedResult: [][]int{
				{1, 2, 2, 2},
				{1, 2, 2, 2},
				{1, 4, 4, 4},
				{3, 4, 4, 4},
			}},
			{topLeftX: 0, topLeftY: 1, width: 4, height: 4, expectedResult: [][]int{
				{1, 1, 1, 2},
				{1, 1, 1, 4},
				{3, 3, 3, 4},
				{3, 3, 3, 4},
			}},
			{topLeftX: 1, topLeftY: 1, width: 4, height: 4, expectedResult: [][]int{
				{1, 1, 2, 2},
				{1, 1, 4, 4},
				{3, 3, 4, 4},
				{3, 3, 4, 4},
			}},
			{topLeftX: 2, topLeftY: 1, width: 4, height: 4, expectedResult: [][]int{
				{1, 2, 2, 2},
				{1, 4, 4, 4},
				{3, 4, 4, 4},
				{3, 4, 4, 4},
			}},
			{topLeftX: 4, topLeftY: 2, width: 4, height: 4, expectedResult: [][]int{
				{4, 4, -1, -1},
				{4, 4, -1, -1},
				{4, 4, -1, -1},
				{-1, -1, -1, -1},
			}},
			{topLeftX: -2, topLeftY: -1, width: 4, height: 4, expectedResult: [][]int{
				{-1, -1, -1, -1},
				{-1, -1, 1, 1},
				{-1, -1, 1, 1},
				{-1, -1, 1, 1},
			}},
			{topLeftX: -1, topLeftY: 3, width: 8, height: 3, expectedResult: [][]int{
				{-1, 3, 3, 3, 4, 4, 4, -1},
				{-1, 3, 3, 3, 4, 4, 4, -1},
				{-1, -1, -1, -1, -1, -1, -1, -1},
			}},
		},
	}

	for i, rSet := range requests {

		for _, r := range rSet {

			res := make([][]int, r.height)
			for y := range res {
				res[y] = make([]int, r.width)
			}

			trees[i].GetContent(r.topLeftX, r.topLeftY, res)

			if !equalContent(r.expectedResult, res) {
				treeAsText := trees[i].GetAsText("fullContent")
				t.Fatalf("Le résultat attendu pour le quadtree \n%s\n avec (topLeftX, topLeftY) = (%d, %d) était le tableau %v mais le résultat retourné par GetContent est %v", treeAsText, r.topLeftX, r.topLeftY, r.expectedResult, res)
			}
		}

	}

}

// Vérifie si deux tableaux de tableaux d'entiers sont
// égaux ou pas.
func equalContent(c1, c2 [][]int) bool {

	if len(c1) != len(c2) {
		return false
	}

	for y := 0; y < len(c1); y++ {

		if len(c1[y]) != len(c2[y]) {
			return false
		}

		for x := 0; x < len(c1[y]); x++ {
			if c1[y][x] != c2[y][x] {
				return false
			}
		}

	}

	return true

}
