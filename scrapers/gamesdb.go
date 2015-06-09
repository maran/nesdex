package scrapers

import (
	"encoding/xml"
	"fmt"
	"github.com/maran/nesdex/common"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Scraper interface {
	UpdateRomFromApi(*common.Rom)
	foreignKey(string)
}

type GamesDb struct {
	Scraper
}

type Data struct {
	GameList []Game `xml:"Game"`
}

type Game struct {
	Id        int      `xml:"id"`
	GameTitle string   `xml:"GameTitle"`
	Images    []Images `xml:"Images"`
}

type Images struct {
	BoxArt []BoxArt `xml:"boxart"`
}

type BoxArt struct {
	Side   string `xml:"side,attr"`
	Thumb  string `xml:"thumb,attr"`
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
	Src    string `xml:",chardata"`
}

func (self *GamesDb) foreignKey() string {
	return "GamesDb"
}

func (self *GamesDb) UpdateRomFromApi(rom *common.Rom) *common.Rom {
	name := url.QueryEscape(strings.Replace(rom.GoodName, "-", ":", -1))
	platform := url.QueryEscape("Nintendo Entertainment System (NES)")
	url := fmt.Sprintf("http://thegamesdb.net/api/GetGame.php?exactname=%s&platform=%s", name, platform)
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
		d := Data{}
		err = xml.Unmarshal(contents, &d)
		if err != nil {
			fmt.Println("Error parsing XML:", err)
			// TODO: Add missing data flag
		}
		if len(d.GameList) > 0 {
			log.Println("Found hit for Rom on GamesDB:", d)
			gameInfo := d.GameList[0]
			rom.GamesDbId = gameInfo.Id
			if len(gameInfo.Images) > 0 {
				for _, img := range gameInfo.Images {
					for _, box := range img.BoxArt {
						src := fmt.Sprintf("http://thegamesdb.net/banners/%s", box.Src)
						rom.BoxArts = append(rom.BoxArts, common.BoxArt{box.Height, box.Width, src, box.Side})
					}
				}
			}
		} else {
			log.Println("No hit found on GamesDB")
		}
	}
	return rom
}
func test() {
	xmlFile, err := os.Open("test.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	//xmlData, err := ioutil.ReadFile("test.xml")
	//	log.Panic("AD")
}
