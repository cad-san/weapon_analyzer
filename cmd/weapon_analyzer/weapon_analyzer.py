#!/usr/bin/env python
# coding:utf-8

import argparse
import sqlite3
import community
import networkx as nx
import matplotlib.pyplot as plt

def getWeaponMain(weapon):
    dict = {
        "52gal":                "52gal",
        "52gal_deco":           "52gal",
        "96gal":                "96gal",
        "96gal_deco":           "96gal",
        "bold":                 "bold",
        "bold_neo":             "bold",
        "dualsweeper":          "dualsweeper",
        "dualsweeper_custom":   "dualsweeper",
        "h3reelgun":            "h3reelgun",
        "h3reelgun_d":          "h3reelgun",
        "heroshooter_replica":  "sshooter",
        "hotblaster":           "hotblaster",
        "hotblaster_custom":    "hotblaster",
        "jetsweeper":           "jetsweeper",
        "jetsweeper_custom":    "jetsweeper",
        "l3reelgun":            "l3reelgun",
        "l3reelgun_d":          "l3reelgun",
        "longblaster":          "longblaster",
        "longblaster_custom":   "longblaster",
        "momiji":               "wakaba",
        "nova":                 "nova",
        "nova_neo":             "nova",
        "nzap85":               "nzap",
        "nzap89":               "nzap",
        "octoshooter_replica":  "sshooter",
        "prime":                "prime",
        "prime_collabo":        "prime",
        "prime_berry":          "prime",
        "promodeler_mg":        "promodeler",
        "promodeler_rg":        "promodeler",
        "rapid":                "rapid",
        "rapid_deco":           "rapid",
        "rapid_elite":          "rapid_elite",
        "rapid_elite_deco":     "rapid_elite",
        "sharp":                "sharp",
        "sharp_neo":            "sharp",
        "sshooter":             "sshooter",
        "sshooter_collabo":     "sshooter",
        "sshooter_wasabi":      "sshooter",
        "wakaba":               "wakaba",
        "carbon":               "carbon",
        "carbon_deco":          "carbon",
        "dynamo":               "dynamo",
        "dynamo_tesla":         "dynamo",
        "dynamo_burned":        "dynamo",
        "heroroller_replica":   "splatroller",
        "hokusai":              "hokusai",
        "hokusai_hue":          "hokusai",
        "pablo":                "pablo",
        "pablo_hue":            "pablo",
        "pablo_permanent":      "pablo",
        "splatroller":          "splatroller",
        "splatroller_collabo":  "splatroller",
        "bamboo14mk1":          "bamboo14",
        "bamboo14mk2":          "bamboo14",
        "bamboo14mk3":          "bamboo14",
        "herocharger_replica":  "splatcharger",
        "liter3k":              "liter3k",
        "liter3k_custom":       "liter3k",
        "liter3k_scope":        "liter3k_scope",
        "liter3k_scope_custom": "liter3k_scope",
        "splatcharger":         "splatcharger",
        "splatcharger_wakame":  "splatcharger",
        "splatscope":           "splatscope",
        "splatscope_wakame":    "splatscope",
        "squiclean_a":          "squiclean",
        "squiclean_b":          "squiclean",
        "squiclean_g":          "squiclean",
        "bucketslosher":        "bucketslosher",
        "bucketslosher_deco":   "bucketslosher",
        "bucketslosher_soda":   "bucketslosher",
        "hissen":               "hissen",
        "hissen_hue":           "hissen",
        "screwslosher":         "screwslosher",
        "screwslosher_neo":     "screwslosher",
        "barrelspinner":        "barrelspinner",
        "barrelspinner_deco":   "barrelspinner",
        "hydra":                "hydra",
        "hydra_custom":         "hydra",
        "splatspinner":         "splatspinner",
        "splatspinner_collabo": "splatspinner",
        "splatspinner_repair":  "splatspinner"
    }
    return dict[weapon]

def getBattle(cur, nextOf):
    battleQuery = u"select battle_id from battle where battle_id > ? order by battle_id ASC limit ?"
    cur.execute(battleQuery, (nextOf, 1))
    cols = cur.fetchall()
    if len(cols) == 0:
        return (None)

    id = 0
    for row in cols:
        id = row[0]
    teamQuery = u"select weapon from team where team_type == ? and battle_id == ?"

    teamA = []
    for row in cur.execute(teamQuery, ("A", id)):
        weapon = row[0]
        if weapon == "":
            continue
        teamA.append(getWeaponMain(weapon))
    teamB = []
    for row in cur.execute(teamQuery, ("B", id)):
        weapon = row[0]
        if weapon == "":
            continue
        teamB.append(getWeaponMain(weapon))
    return (id, teamA, teamB)

def updateGraph(G, teamA, teamB):
    for weapon in teamA:
        G.add_node(weapon)
    for weapon in teamB:
        G.add_node(weapon)

    for memberA in teamA:
        for memberB in teamB:
            if G.has_edge(memberA,memberB):
                G[memberA][memberB]['weight'] += 1
            else:
                G.add_edge(memberA,memberB,weight=1)

if __name__ == "__main__":

    parser = argparse.ArgumentParser(description='analyze weapon from database')
    parser.add_argument('database')
    args = parser.parse_args()

    db = sqlite3.connect(args.database)
    cur = db.cursor()

    G = nx.Graph()
    id = 0
    while(1):
        res = getBattle(cur, id)
        if res is None:
            break
        (id, teamA, teamB) = res
        updateGraph(G,teamA,teamB)

    db.close()

    for u, v, d in G.edges(data=True):
        if d['weight'] == 1:
            G.remove_edge(u,v)
        print((u,v),d['weight'])

    part = community.best_partition(G)
    values = [part.get(node) for node in G.nodes()]

    pos = nx.spring_layout(G)
    nx.draw_networkx_nodes(G, pos, node_size = 200,
                                node_color =values,
                                cmap = plt.get_cmap('jet'))
    nx.draw_networkx_labels(G, pos)
    plt.show()
