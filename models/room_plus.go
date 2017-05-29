package models

import (
	"fmt"
	"sort"
)

type HuiZong struct {
	Master   string
	NicName  string
	ScoreSum int
	LenUsers int
	Users    []Player
}

type sortPlayer []Player

func (su sortPlayer) Less(i, j int) bool {
	return su[i].Score > su[j].Score
}

func (su sortPlayer) Len() int {
	return len(su)
}

func (su sortPlayer) Swap(i, j int) {
	su[i], su[j] = su[j], su[i]
}

func (r *Room) Hui(master string) (HuiZong, error) {

	huizong := HuiZong{
		Master:   master,
		LenUsers: len(r.users),
	}

	pl, ok := r.users[master]
	if !ok {
		return huizong, fmt.Errorf("no master in the room")
	}
	huizong.NicName = pl.NicName

	for key := range r.users {

		if r.users[key].Score != 0 {
			huizong.ScoreSum += r.users[key].Score
			huizong.Users = append(huizong.Users, *(r.users[key]))
		}
	}

	sort.Sort(sortPlayer(huizong.Users))
	return huizong, nil
}
