package main

import (
	"github.com/maran/nesdex/common"
	"github.com/maran/nesdex/scrapers"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

func main() {
	db := common.NewDatabase()

	scraper := scrapers.GamesDb{}

	var romGroups []common.RomGroup
	err := db.RomGroupCollection.Find(bson.M{"initial_games_db_scan": false}).All(&romGroups)
	if err != nil {
		log.Panic("Error fetching roms", err)
	}
	log.Println("Found", len(romGroups), "rom groups to check on GamesDB")
	for _, romGroup := range romGroups {
		log.Println("Scraping info for", romGroup.Name)
		scraper.UpdateRomFromApi(&romGroup)
		romGroup.GamesDbScanned = true
		db.PersistRomGroup(&romGroup)
		time.Sleep(500 * time.Millisecond)
	}
	log.Println("Done")
	err = db.RomGroupCollection.Find(bson.M{"games_db_id": 0, "initial_games_db_scan": true}).All(&romGroups)
	if err != nil {
		log.Panic("Error fetching roms", err)
	}
	log.Println("Found", len(romGroups), "roms that were checked but not found. Rechecking")
	for _, romGroup := range romGroups {
		log.Println("Name:", romGroup.Name)
		scraper.UpdateRomFromApi(&romGroup)
		db.PersistRomGroup(&romGroup)
		time.Sleep(500 * time.Millisecond)
	}
}
