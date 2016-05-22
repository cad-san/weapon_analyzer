package ikaWeapon

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/naoina/genmai"
)

type BattleTable struct {
	BattleID int    `column:"battle_id" db:"unique"`
	Lobby    string `column:"lobby`
	Rule     string `column:"rule"`
}

type TeamTable struct {
	ID       int64  `db:"pk"`
	BattleID int    `column:"battle_id"`
	TeamType string `column:"team_type"`
	Weapon   string `column:"weapon"`
}

type btSlice []BattleTable

type DBHandle struct {
	genmai.DB
}

func (b *BattleTable) TableName() string {
	return "battle"
}

func (t *TeamTable) TableName() string {
	return "team"
}

// CreateDBHandle return DBHandler which provide DB Access
func CreateDBHandle(path string) (*DBHandle, error) {
	db, err := genmai.New(&genmai.SQLite3Dialect{}, path)
	return &DBHandle{DB: *db}, err
}

func (db *DBHandle) PrepareBattle() error {
	err := db.CreateTableIfNotExists(&BattleTable{})
	if err != nil {
		return err
	}
	return db.CreateTableIfNotExists(&TeamTable{})
}

func (db *DBHandle) InsertBattle(b *Battle) (err error) {
	battle, team := b.convertColumns()
	defer func() {
		if e := recover(); e != nil {
			db.Rollback()
			err = fmt.Errorf("insert battle failed:%v", e)
		} else {
			db.Commit()
			err = nil
		}
	}()
	if err = db.Begin(); err != nil {
		panic(err)
	}
	if _, err = db.Insert(&battle); err != nil {
		panic(err)
	}
	if _, err = db.Insert(&team); err != nil {
		panic(err)
	}
	return nil
}

func (db *DBHandle) GetBattles(nextOf int, num int) ([]Battle, error) {
	var bTables []BattleTable
	err := db.Select(&bTables, db.From(&BattleTable{}), db.Where("battle_id", ">", nextOf), db.OrderBy("battle_id", genmai.ASC).Limit(num))
	if err != nil {
		return nil, err
	}
	idList := btSlice(bTables).getIDs()
	var teams []TeamTable
	err = db.Select(&teams, db.From(&TeamTable{}), db.Where("battle_id").In(idList))
	if err != nil {
		return nil, err
	}
	var battle []Battle
	for _, bt := range bTables {
		battle = append(battle, bt.convert(teams))
	}
	return battle, err
}

func (b *Battle) IsValid() bool {
	if b.Lobby != "standard" {
		return false
	}
	if b.Rule == "nawabari" {
		return false
	}
	if b.TeamA.Weapon == nil || len(b.TeamA.Weapon) == 0 {
		return false
	}
	if b.TeamB.Weapon == nil || len(b.TeamB.Weapon) == 0 {
		return false
	}
	return true
}

func (b *Battle) convertColumns() (BattleTable, []TeamTable) {
	bTable := BattleTable{BattleID: b.ID, Lobby: b.Lobby, Rule: b.Rule}
	var team []TeamTable
	for _, w := range b.TeamA.Weapon {
		column := TeamTable{
			BattleID: b.ID,
			TeamType: "A",
			Weapon:   w,
		}
		team = append(team, column)
	}
	for _, w := range b.TeamB.Weapon {
		column := TeamTable{
			BattleID: b.ID,
			TeamType: "B",
			Weapon:   w,
		}
		team = append(team, column)
	}
	return bTable, team
}

func (b *BattleTable) convert(team []TeamTable) Battle {
	battle := Battle{ID: b.BattleID, Lobby: b.Lobby, Rule: b.Rule}
	for _, t := range team {
		if t.BattleID != b.BattleID {
			continue
		}
		switch t.TeamType {
		case "A":
			battle.TeamA.Weapon = append(battle.TeamA.Weapon, t.Weapon)
		case "B":
			battle.TeamB.Weapon = append(battle.TeamB.Weapon, t.Weapon)
		}
	}
	return battle
}

func (bts btSlice) getIDs() []int {
	slice := []BattleTable(bts)
	var ids []int
	for _, b := range slice {
		ids = append(ids, b.BattleID)
	}
	return ids
}
