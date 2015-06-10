package main

import (
	"flag"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/jinzhu/gorm"
	"github.com/maran/nesdex/common"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

var romFolder = flag.String("romFolder", "roms", "Path to your rom folder")
var configPath = flag.String("configPath", common.GetUserDir()+"/.config/goodfrontend", "Path to your config")
var apiLocation = flag.String("apiLocation", "127.0.0.1:8080", "Path to scrape information from")

func main() {
	flag.Parse()
	frontEnd := newGoodFront(*configPath)
	scanner := NewScanner(frontEnd)
	scanner.ScanFolder()

	api := Api{frontEnd, "localhost:8888"}
	api.Start()
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

type Api struct {
	front      *GoodFront
	listenPort string
}
type RomsResponse struct {
	Roms []RomFile `json:"roms"`
	Page int       `json:"page"`
}

func (self *Api) Start() {
	newApi := rest.NewApi()
	newApi.Use(rest.DefaultDevStack...)
	newApi.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	router, err := rest.MakeRouter(
		&rest.Route{"GET", "/roms/:page", self.getRoms},
		&rest.Route{"GET", "/roms/filter/:term", self.getRoms},
		&rest.Route{"GET", "/start/:id", self.startRom},
	)

	if err != nil {
		log.Fatal(err)
	}

	newApi.SetApp(router)
	log.Println("Server started and listening on %s", self.listenPort)
	log.Fatal(http.ListenAndServe(self.listenPort, newApi.MakeHandler()))
}
func (self *Api) startRom(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	log.Println("Starting rom with ID", id)

	romId, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		return
	}

	rom := RomFile{ID: int(romId)}
	self.front.db.Find(&rom)

	cmd := exec.Command("nes", rom.FullPath)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}
func (self *Api) getRoms(w rest.ResponseWriter, r *rest.Request) {
	page := r.PathParam("page")
	perPage := int(100)
	var roms []RomFile
	pageInt, err := strconv.ParseInt(page, 10, 32)
	if err != nil {
		log.Println("Error could not parse page assuming 1")
		pageInt = 1
		return
	}

	offset := int(pageInt) * perPage
	self.front.db.Offset(offset).Limit(perPage).Preload("BoxArt").Find(&roms) //.RomCollection.Find(bson.M{"verified": true, "initial_games_db_scan": true}).Skip(offset).Limit(perPage).All(&roms)
	if err != nil {
		log.Println("Error fetching roms", err)
	}
	w.WriteJson(RomsResponse{Roms: roms, Page: int(pageInt)})
}
