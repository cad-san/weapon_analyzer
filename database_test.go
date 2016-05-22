package ikaWeapon

import (
	"testing"
)

var dummyBattle = Battle{
	ID:    10000,
	Lobby: "standard",
	Rule:  "area",
	TeamA: Team{
		Weapon: []string{
			"96gal_deco",
			"liter3k_scope",
			"splatspinner_collabo",
			"hydra_custom",
		},
	},
	TeamB: Team{
		Weapon: []string{
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

func TestGetBattle(t *testing.T) {
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
	b, err := handle.GetBattles(0, 1)
	if err != nil {
		t.Fatalf("get error :%v", err)
	}
	if len(b) != 1 {
		t.Errorf("length of battle slice must be 1 but%v", len(b))
	}
}
