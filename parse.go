package ikaWeapon

import (
	"encoding/json"
	"io"
)

type strmap struct {
	Key string
}

type player struct {
	Team   string
	Weapon strmap
}

type rawPlayerJSON struct {
	ID      int
	Lobby   strmap
	Rule    strmap
	Players []player
}

func decodeJSONBattleSingle(r io.Reader) (*Battle, error) {
	d := json.NewDecoder(r)
	raw := rawPlayerJSON{}

	err := d.Decode(&raw)
	if err != nil {
		return nil, err
	}
	battle := raw.getBattle()
	return &battle, nil
}

func decodeJSONBattle(r io.Reader) ([]Battle, error) {
	d := json.NewDecoder(r)
	raw := []rawPlayerJSON{}

	err := d.Decode(&raw)
	if err != nil {
		return nil, err
	}

	var battle []Battle
	for _, r := range raw {
		battle = append(battle, r.getBattle())
	}
	return battle, nil
}

func (r rawPlayerJSON) getBattle() Battle {
	b := Battle{ID: r.ID, Lobby: r.Lobby.Key, Rule: r.Rule.Key}
	for _, p := range r.Players {
		switch p.Team {
		case "my":
			b.TeamA.Weapon = append(b.TeamA.Weapon, p.Weapon.Key)
		case "his":
			b.TeamB.Weapon = append(b.TeamB.Weapon, p.Weapon.Key)
		default:
			// do nothing
		}
	}
	return b
}
