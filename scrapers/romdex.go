package scrapers

import (
	"encoding/json"
	"fmt"
	"github.com/maran/nesdex/common"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type RomDex struct {
	Scraper
}

const RomDexHost = "http://localhost:8080"

func (self *RomDex) foreignKey() string {
	return "RomDex"
}

func (self *RomDex) UpdateRomFromApi(rom *common.Rom) *common.RomResponse {
	url := fmt.Sprintf("%s/api/v1/roms/md5/%s", RomDexHost, rom.Md5)
	log.Println("Grabbing URL:", url)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Got HTTP error during scrape %s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		log.Println(string(contents))
		d := common.RomResponse{}
		err = json.Unmarshal(contents, &d)
		if err != nil {
			log.Println("Error from RomDex")
		}
		return &d
	}
	/** Goodfrontend rom response changes **/
	return &common.RomResponse{Rom: *rom}
}
