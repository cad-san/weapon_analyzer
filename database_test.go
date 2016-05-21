package ikaWeapon

import (
	"testing"
)

var dummyBattle = Battle{
	id:   10000,
	rule: "area",
	teamA: Team{
		weapon: []string{
			"96gal_deco",
			"liter3k_scope",
			"splatspinner_collabo",
			"hydra_custom",
		},
	},
	teamB: Team{
		weapon: []string{
			"sshooter_collabo",
			"wakaba",
			"52gal",
			"liter3k_scope_custom",
		},
	},
}

func TestInsertBattle(t *testing.T) {
	handle, err := CreateDBHandle(":memory:")
	if err != nil {
		t.Fatalf("create database error :%v", err)
	}
	defer handle.Close()
	err = handle.PrepareBattle()
	if err != nil {
		t.Fatalf("create table error :%v", err)
	}
	err = handle.InsertBattle(&dummyBattle)
	if err != nil {
		t.Fatalf("insert error :%v", err)
	}
}
