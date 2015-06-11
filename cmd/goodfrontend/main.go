package main

import (
	"flag"
	"github.com/jinzhu/gorm"
	"github.com/maran/nesdex/common"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

var romFolder = flag.String("romFolder", "roms", "Path to your rom folder")
var configPath = flag.String("configPath", common.GetUserDir()+"/.config/goodfrontend", "Path to your config")
var apiLocation = flag.String("apiLocation", "127.0.0.1:8888", "API Location to scrape information from")

func main() {
	flag.Parse()

	frontEnd := newGoodFront(*configPath)

	scanner := NewScanner(frontEnd)
	scanner.ScanFolder()

	// Setup ROM API
	api := Api{frontEnd, *apiLocation}
	go api.Start()

	// Server HTML/JS
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./../../gui"))))
	log.Fatal(http.ListenAndServe(":8181", nil))
}

type GoodFront struct {
	db          gorm.DB
	romFolder   string
	configPath  string
	apiLocation string
}

func newGoodFront(dbDir string) *GoodFront {
	dbPath := *configPath + "/goodfrontend.db"

	log.Println("Loading dbPath", dbPath)
	err := os.MkdirAll(dbDir, 0700) // read, write and dir search for user
	if err != nil {
		log.Fatal("Error creating database folder", err)
	}

	frontEnd := GoodFront{romFolder: *romFolder, configPath: *configPath, apiLocation: *apiLocation}
	frontEnd.db, err = gorm.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}

	frontEnd.db.AutoMigrate(&RomFile{})
	frontEnd.db.AutoMigrate(&BoxArtDetail{})
	frontEnd.db.LogMode(true)
	return &frontEnd
}
