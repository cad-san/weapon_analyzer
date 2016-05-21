package ikaWeapon

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/naoina/genmai"
)

type BattleTable struct {
	BattleID int    `column:"battle_id"`
	Rule     string `column:"rule"`
}

type TeamTable struct {
	ID       int64  `db:"pk"`
	BattleID int    `column:"battle_id"`
	TeamType string `column:"team_type"`
	Weapon   string `column:"weapon"`
}

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

func (b *Battle) convertColumns() (BattleTable, []TeamTable) {
	bTable := BattleTable{BattleID: b.id, Rule: b.rule}
	var team []TeamTable
	for _, w := range b.teamA.weapon {
		column := TeamTable{
			BattleID: b.id,
			TeamType: "A",
			Weapon:   w,
		}
		team = append(team, column)
	}
	for _, w := range b.teamB.weapon {
		column := TeamTable{
			BattleID: b.id,
			TeamType: "B",
			Weapon:   w,
		}
		team = append(team, column)
	}
	return bTable, team
}
