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
	log.Println("Found", len(roms), "roms to check on GamesDB")
	for _, rom := range roms {
		log.Println("Scraping info for", rom.GoodName)
		scraper.UpdateRomFromApi(&rom)
		rom.GamesDbScanned = true
		db.PersistRom(&rom)
		time.Sleep(500 * time.Millisecond)
	}
	log.Println("Done")
	err = db.RomCollection.Find(bson.M{"games_db_id": 0, "initial_games_db_scan": true}).All(&roms)
	if err != nil {
		log.Panic("Error fetching roms", err)
	}
	log.Println("Found", len(roms), "roms that were checked but not found. Rechecking")
	for _, rom := range roms {
		log.Println("Name:", rom.GoodName)
		scraper.UpdateRomFromApi(&rom)
		db.PersistRom(&rom)
		time.Sleep(500 * time.Millisecond)
	}
	// TODO: Build a system that either merges the same rom in different countries to share found data
}
