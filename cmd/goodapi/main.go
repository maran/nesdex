package main

import (
	"flag"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/maran/nesdex/common"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strconv"
)

var listenPort = flag.String("listenPort", "0.0.0.0:8080", "Port where the json api should listen at in host:port format.")

type NesApi struct {
	db *common.Database
}

type RomsResponse struct {
	Roms []common.Rom `json:"roms"`
	Page int          `json:"page"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func (self *NesApi) getRoms(w rest.ResponseWriter, r *rest.Request) {
	page := r.PathParam("page")
	perPage := int(100)
	var roms []common.Rom
	pageInt, err := strconv.ParseInt(page, 10, 32)
	if err != nil {
		log.Println("Error could not parse page assuming 1")
		pageInt = 1
		return
	}

	// Can we make this DRYer?
	offset := int(pageInt) * perPage
	err = self.db.RomCollection.Find(bson.M{"verified": true, "initial_games_db_scan": true}).Skip(offset).Limit(perPage).All(&roms)
	if err != nil {
		log.Println("Error fetching roms", err)
	}
	w.WriteJson(RomsResponse{Roms: roms, Page: int(pageInt)})
}

func (self *NesApi) getRom(w rest.ResponseWriter, r *rest.Request) {
	md5 := r.PathParam("md5")
	var rom common.Rom
	var romGroup common.RomGroup

	err := self.db.RomCollection.Find(bson.M{"md5": md5}).One(&rom)
	if err != nil {
		log.Println("Error fetching Rom with md5", md5)
		w.WriteJson(ErrorResponse{404, true, "Rom not found in database"})
		return
	}
	log.Println("Looking for", rom.RomGroupId)
	err = self.db.RomGroupCollection.Find(bson.M{"_id": rom.RomGroupId}).One(&romGroup)
	if err != nil {
		log.Println("No rom collection found", err)
	}
	w.WriteJson(common.RomResponse{Rom: rom, RomGroup: romGroup})
}

func main() {
	flag.Parse()

	db := common.NewDatabase()
	nesApi := NesApi{db: db}
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
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
		&rest.Route{"GET", "/roms/:page", nesApi.getRoms},
		&rest.Route{"GET", "/roms/filter/:term", nesApi.getRoms},
		&rest.Route{"GET", "/roms/md5/:md5", nesApi.getRom},
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	http.Handle("/api/v1/", http.StripPrefix("/api/v1", api.MakeHandler()))
	log.Println("Server started and listening on %s", *listenPort)
	log.Fatal(http.ListenAndServe(*listenPort, nil))

}
