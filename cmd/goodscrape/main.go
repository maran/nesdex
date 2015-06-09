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

	var roms []common.Rom
	err := db.RomCollection.Find(bson.M{"verified": true, "initial_games_db_scan": false}).All(&roms)
	if err != nil {
		log.Panic("Error fetching roms", err)
	}
	for _, rom := range roms {
		log.Println("Scraping info for", rom.GoodName)
		scraper.UpdateRomFromApi(&rom)
		rom.GamesDbScanned = true
		db.PersistRom(&rom)
		time.Sleep(1000 * time.Millisecond)
	}
}
