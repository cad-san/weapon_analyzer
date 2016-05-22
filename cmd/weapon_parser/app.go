package main

import (
	"flag"
	"fmt"
	"github.com/cad-san/weapon_analyzer"
)

var (
	startID     int
	dbPath      string
	numRequest  int
	totalBattle int
)

func initFlags() {
	flag.StringVar(&dbPath, "output", "test.db", "file path of SQLite database")
	flag.StringVar(&dbPath, "o", "test.db", "file path of SQLite database")
	flag.IntVar(&startID, "start", 0, "get id")
	flag.IntVar(&startID, "s", 0, "get id")
	flag.IntVar(&numRequest, "request", 10, "number of battles par http request")
	flag.IntVar(&numRequest, "r", 10, "number of battles par http request")
	flag.IntVar(&totalBattle, "total", 10, "number of total battles")
	flag.IntVar(&totalBattle, "t", 10, "number of total battles")
	flag.Parse()
}

func parseAndSave() error {
	db, err := ikaWeapon.CreateDBHandle(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.PrepareBattle()
	if err != nil {
		return err
	}
	client := ikaWeapon.CreateClient()

	count := 0
	nextID := startID
	for count < totalBattle {
		fmt.Printf("get request : %d\n", nextID)
		list, err := client.GetBattleOlderThan(nextID, numRequest)
		if err != nil {
			return err
		}

		for _, b := range list {
			if b.IsValid() {
				if err = db.InsertBattle(&b); err != nil {
					fmt.Println(err)
					continue
				}
				count++
				fmt.Println(b)
			}
			if count >= totalBattle {
				break
			}
		}
		nextID = list[len(list)-1].ID
	}
	return nil
}

func main() {
	initFlags()
	if err := parseAndSave(); err != nil {
		panic(err)
	}
}
