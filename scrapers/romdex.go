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
	Host string
}

func (self *RomDex) foreignKey() string {
	return "RomDex"
}

func (self *RomDex) UpdateRomFromApi(rom *common.Rom) *common.Rom {
	url := fmt.Sprintf("http://localhost:8080/roms/md5/%s", rom.Md5)
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
		d := common.Rom{}
		err = json.Unmarshal(contents, &d)
		if err != nil {
			log.Println("Error from RomDex")
		}
		return &d
	}
	return rom
}
