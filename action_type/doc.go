package action_type

const (
	PlayerSpawnId = iota
	PlayerDespawnId
	PlayerMoveId
	OrientationChangeId
	AskChunkActionId
	SetChunkActionId
)

type PlayerInitData struct {
	Name              string
	X, Y, Orientation int
	PlaceBlock        []int
}

type InitWorldState struct {
	// Pour l’instant, tout ce que les peer doivent savoir sur la connexion est la position du joueur
	Players []PlayerInitData
}
