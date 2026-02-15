package action_out

import "gitlab.univ-nantes.fr/jezequel-l/quadtree/action_type"

type NotifyPlayerDepawn struct {
	To_  string
	Name string
}

func NewNotifyPlayerDespawn(to, name string) NotifyPlayerDepawn {
	return NotifyPlayerDepawn{
		To_:  to,
		Name: name,
	}
}

func (a NotifyPlayerDepawn) GetData() []byte {
	return append([]byte{action_type.PlayerDespawnId}, []byte(a.Name)...)
}

func (a NotifyPlayerDepawn) To() string {
	return a.To_
}
