package game

import (
	//Ajouté par moi
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/camera"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/character"
	"gitlab.univ-nantes.fr/jezequel-l/quadtree/floor"
)

type ActionOut interface {
	GetData() []byte
	To() string
}

type GameAction interface {
	Execute(g *Game)
}

// Game est le type permettant de représenter les données du jeu.
// Aucun champs n'est exporté pour le moment.
//
// Les champs non exportés sont :
//   - camera : la représentation de la caméra
//   - floor : la représentation du terrain
//   - character : la représentation du personnage
type Game struct {
	//Pour le bon fonctionnement du jeu
	Camera    camera.Camera
	Floor     floor.Floor
	Character character.Character

	// Variables de fonctionnement
	PlayerOnUi string // Ui actuellement affichée au joueur parmi ( "" > rien | "inventory" > inventaire | "crafting_table" > table de fabrication | "furnace" > four ) et peut être plus tard "menu"

	// Multijoueur
	OtherPlayers    []character.Character
	ActionOutChan   chan ActionOut
	GameActionChan  chan GameAction
	ChunkLoadedChan chan ChunkLoaded

	PlaceBlocX  int
	PlaceBlocY  int
	PlaceBloc   bool
	PlaceBlocID int
}
