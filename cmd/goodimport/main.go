package main

import (
	"flag"
	"fmt"
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

	fmt.Println("Scanning", *romsPath)
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
		fmt.Println("Checking", b)
		name := fullNameMatch.FindString(b)

		// Remove trailer whitespace
		name = name[0 : len(name)-1]
		rom := common.NewRom(name)
		//searchApi(name)
		log.Println("Found", name, "variation")
		rom.Md5 = common.CalcMd5(filePath)
		rom.Name = b

		result := trnMatch.FindAllStringSubmatch(b, -1)
		for _, res := range result {
			fmt.Println("Trainer:", res[3])
			rom.Trainer = true
			rom.TrainerName = res[3]
		}
		result = parMatch.FindAllStringSubmatch(b, -1)
		for _, res := range result {
			if val, ok := common.CountriesMap[res[1]]; ok {
				rom.CountryId = val
			}

			switch string(res[1]) {
			case "Hack":
				fmt.Println("Unspecified hack")
				rom.Hack = true
			// Revisions
			case "PRG0":
				fmt.Println("Program Revision 0")
			case "PRG1":
				fmt.Println("Program Revision 1")
			default:
				found := false
				// Check if the result includes the word Hack
				if strings.Contains(res[1], "Hack") {
					fmt.Println("Custom hack:", res[1])
					rom.Hack = true
					rom.HackName = res[1]
					found = true
				}
				// Check if the word Mapper is included
				var mprMatch = regexp.MustCompile(`Mapper (\d*)`)
				p := mprMatch.FindString(res[1])
				if len(p) > 0 {
					fmt.Println("Mapper:", p)
					found = true
				}
				if found == false {
					fmt.Println("Unprocessed:", res[1])
				}
			}
		}
		result = stdMatch.FindAllStringSubmatch(b, -1)
		for _, res := range result {
			fmt.Println(result)
			switch string(res[1][0]) {
			case "!":
				fmt.Println("Verified")
				rom.Verified = true
			case "a":
				fmt.Println("Alternative version")
				rom.Alternative = true
			case "b":
				fmt.Println("Bad dump")
				rom.BadDump = true
			case "f":
				fmt.Println("Fixed")
				rom.Fixed = true
			case "h":
				fmt.Println("Hacked")
				rom.Hack = true
			case "o":
				fmt.Println("Overdump")
				rom.Overdump = true
			case "p":
				fmt.Println("Pirated")
				rom.Pirated = true
			case "t":
				fmt.Println("Trainer")
				rom.Trainer = true
			default:
				if res[1] == "!p" {
					fmt.Println("Best but waiting for dump")
				} else {
					fmt.Println("Unknown code", res[1])
				}
			}

		}
		fmt.Println("---------------------------------")
		db.PersistRom(rom)
	}
}
