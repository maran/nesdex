package main

import (
	"github.com/maran/nesdex/common"
	"github.com/maran/nesdex/scrapers"
	"log"
)

type RomFile struct {
	common.Rom
	ID         int
	FullPath   string `json:"full_path", bson:"-"`
	Identified bool   `json:"identified"`
	BoxArt     []BoxArtDetail
}

type BoxArtDetail struct {
	common.BoxArt
	ID        int
	RomFileId int
}

func (self *RomFile) Identify() bool {
	dex := scrapers.RomDex{}
	log.Println("Scraping info for", self.Rom.Md5)
	updatedRom := dex.UpdateRomFromApi(&self.Rom)
	self.Rom = *updatedRom
	log.Println("New info", self.Rom)

	for _, art := range updatedRom.BoxArts {
		self.BoxArt = append(self.BoxArt, BoxArtDetail{BoxArt: art})
	}
	return true
}
