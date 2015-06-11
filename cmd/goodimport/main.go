package main

import (
	"flag"
	"github.com/maran/nesdex/common"
	"log"
	"path/filepath"
	"regexp"
	"strings"
)

// Resources used
// http://en.wikipedia.org/wiki/GoodTools
// http://wiki.thegamesdb.net/index.php/GetGame
// http://thegamesdb.net/api/GetGame.php?exactname=Adventure%20Island&platform=Nintendo%20Entertainment%20System%20(NES)
// TODO: Add tons of flags for more information
var romsPath = flag.String("roms_path", "./roms", "path to the roms to import.")

func main() {
	flag.Parse()

	log.Println("Scanning", *romsPath)
	matches, err := filepath.Glob(*romsPath + "/*.nes")
	db := common.NewDatabase()

	var fullNameMatch = regexp.MustCompile(`^([^\(|^\[]*)`)
	var parMatch = regexp.MustCompile(`\(([^\)]+)\)`)
	var stdMatch = regexp.MustCompile(`\[([a-zA-Z*\d*!]*)\]`)
	var trnMatch = regexp.MustCompile(`\[(T)(-|\+)(\w*\d?\.?\s?\w*)\]`)

	if err != nil {
		log.Fatal("Error", err)
	}

	for _, filePath := range matches {
		_, b := filepath.Split(filePath)
		log.Println("Checking", b)
		name := fullNameMatch.FindString(b)

		// Remove trailer whitespace
		name = name[0 : len(name)-1]

		romGroup := db.FindOrCreateGroup(name)
		log.Println("Using Romgroup", name, "id", romGroup.Id)

		md5 := common.CalcMd5(filePath)
		rom := db.FindOrInitializeRom(md5)
		rom.GoodName = name
		rom.Name = b
		rom.RomGroupId = romGroup.Id

		log.Println("Analyzing", b)

		log.Println(rom)

		result := trnMatch.FindAllStringSubmatch(b, -1)
		for _, res := range result {
			log.Println("Trainer:", res[3])
			rom.Trainer = true
			rom.TrainerName = res[3]
		}
		result = parMatch.FindAllStringSubmatch(b, -1)
		for _, res := range result {
			found := false
			if val, ok := common.CountriesMap[res[1]]; ok {
				log.Println("Found country:", res[1])
				found = true
				rom.CountryId = val
			}

			switch string(res[1]) {
			case "Hack":
				log.Println("Unspecified hack")
				rom.Hack = true
			// Revisions
			case "PRG0":
				log.Println("Program Revision 0")
			case "PRG1":
				log.Println("Program Revision 1")
			default:
				// Check if the result includes the word Hack
				if strings.Contains(res[1], "Hack") {
					log.Println("Custom hack:", res[1])
					rom.Hack = true
					rom.HackName = res[1]
					found = true
				}
				// Check if the word Mapper is included
				var mprMatch = regexp.MustCompile(`Mapper (\d*)`)
				p := mprMatch.FindString(res[1])
				if len(p) > 0 {
					log.Println("Mapper:", p)
					found = true
				}
				if found == false {
					log.Println("Unprocessed:", res[1])
				}
			}
		}
		result = stdMatch.FindAllStringSubmatch(b, -1)
		for _, res := range result {
			log.Println(result)
			switch string(res[1][0]) {
			case "!":
				log.Println("Verified")
				rom.Verified = true
			case "a":
				log.Println("Alternative version")
				rom.Alternative = true
			case "b":
				log.Println("Bad dump")
				rom.BadDump = true
			case "f":
				log.Println("Fixed")
				rom.Fixed = true
			case "h":
				log.Println("Hacked")
				rom.Hack = true
			case "o":
				log.Println("Overdump")
				rom.Overdump = true
			case "p":
				log.Println("Pirated")
				rom.Pirated = true
			case "t":
				log.Println("Trainer")
				rom.Trainer = true
			default:
				if res[1] == "!p" {
					log.Println("Best but waiting for dump")
				} else {
					log.Println("Unknown code", res[1])
				}
			}

		}
		log.Println("---------------------------------")
		db.PersistRom(rom)
	}
}
