package main

import (
	"github.com/maran/nesdex/common"
	"log"
	"path/filepath"
)

type Scanner struct {
	frontEnd *GoodFront
}

func (self *Scanner) ScanFolder() int {
	log.Println("Scanning", self.frontEnd.romFolder)
	matches, err := filepath.Glob(self.frontEnd.romFolder + "/*.nes")
	if err != nil {
		log.Println("Error scanning", err)
		return 0
	}

	for _, filePath := range matches {
		_, fileName := filepath.Split(filePath)
		fullPath, err := filepath.Abs(filePath)
		if err != nil {
			log.Println("Skipping value because of path error")
			continue
		}
		var rom RomFile

		bRom := common.Rom{Md5: common.CalcMd5(filePath)}
		self.frontEnd.db.Where(RomFile{Rom: bRom}).FirstOrInit(&rom)

		if rom.Name != "" {
			log.Println("Rom", rom.Name, "with hash", rom.Md5, "already in database")
		} else {
			log.Println("Found new file", fileName, "with hash", rom.Md5)
			rom.Rom.Name = fileName
			rom.FullPath = fullPath
			// TODO: Do this in a seperate loop somewhere via a go-channel
			rom.Identify()
			self.frontEnd.db.Save(&rom)
		}
	}
	return 0
}

func NewScanner(frontEnd *GoodFront) *Scanner {
	return &Scanner{frontEnd}
}
