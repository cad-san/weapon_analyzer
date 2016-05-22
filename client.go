package ikaWeapon

import (
	"net/http"
	"net/url"
	"strconv"
)

const (
	statInkBattleAPI = "https://stat.ink/api/v1/battle"
)

type InkClient struct {
	http.Client
}

type Team struct {
	Weapon []string
}

type Battle struct {
	ID    int
	Lobby string
	Rule  string
	TeamA Team
	TeamB Team
}

func CreateClient() *InkClient {
	client := &InkClient{}
	return client
}

func (c *InkClient) GetBattle(id int) (*Battle, error) {
	query := buildSingleBattleQuery(id)
	url := statInkBattleAPI + "?" + query.Encode()
	resp, err := c.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decodeJSONBattleSingle(resp.Body)
}

func (c *InkClient) GetRecentBattle(count int) ([]Battle, error) {
	return c.getBattleList("", 0, count)
}

func (c *InkClient) GetBattleOlderThan(id int, count int) ([]Battle, error) {
	return c.getBattleList("older_than", id, count)
}

func (c *InkClient) getBattleList(mode string, id int, count int) ([]Battle, error) {
	query := buildBattleQuery(mode, id, count)
	url := statInkBattleAPI + "?" + query.Encode()
	resp, err := c.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	battle, err := decodeJSONBattle(resp.Body)
	return battle, err
}

func buildBattleQuery(mode string, id int, count int) url.Values {
	query := url.Values{}
	if len(mode) > 0 {
		query.Add(mode, strconv.Itoa(id))
	}
	query.Add("count", strconv.Itoa(count))
	return query
}

func buildSingleBattleQuery(id int) url.Values {
	query := url.Values{}
	query.Add("id", strconv.Itoa(id))
	return query
}
