package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
	"os/exec"
	"strconv"
)

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
