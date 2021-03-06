package common

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

const (
	Australia int = iota
	Asia
	Brazil
	Canada
	China
	Dutch
	Europe
	France
	Germany
	Greece
	HongKong
	Italy
	Japan
	Korea
	Netherlands
	Norway
	Russia
	Spain
	Sweden
	UnitedStates
	UnitedKingdom
	World
	Unlicensed
	PublicDomain
	UnknownCountry
)

var CountriesMap = map[string]int{
	"A":   Australia,
	"As":  Asia,
	"B":   Brazil,
	"C":   Canada,
	"Ch":  China,
	"D":   Dutch,
	"E":   Europe,
	"F":   France,
	"G":   Germany,
	"Gr":  Greece,
	"HK":  HongKong,
	"I":   Italy,
	"J":   Japan,
	"K":   Korea,
	"Nl":  Netherlands,
	"No":  Norway,
	"R":   Russia,
	"S":   Spain,
	"Sw":  Sweden,
	"U":   UnitedStates,
	"Uk":  UnitedKingdom,
	"W":   World,
	"Unl": Unlicensed,
	"PD":  PublicDomain,
	"Unk": UnknownCountry,
}

type RomResponse struct {
	RomGroup RomGroup `json:"rom_group"`
	Rom      Rom      `json:"rom"`
}

type RomGroup struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name   string
	System int8 `json:"system"`

	GamesDbId      int  `bson:"games_db_id" json:"games_db_id"`
	GamesDbScanned bool `bson:"initial_games_db_scan" json:"games_db_scanned"`

	BoxArts []BoxArt `bson:"box_arts" json:"box_arts"`
}

type Rom struct {
	Name     string `json:"name"`
	GoodName string `bson:"good_name" json:"sanitized_name"`
	Md5      string `json:"md5"`
	System   int8   `json:"system"`

	Hack        bool          `json:"hack"`
	HackName    string        `bson:"hack_name" json:"hack_name"`
	Trainer     bool          `json:"trainer"`
	TrainerName string        `bson:"trainer_name" json:"trainer_name"`
	BadDump     bool          `bson:"bad_dump" json:"bad_dump"`
	Overdump    bool          `json:"overdump"`
	Verified    bool          `json:"verified"`
	Alternative bool          `json:"alternative"`
	Pirated     bool          `json:"pirated"`
	Fixed       bool          `json:"fixed"`
	CountryId   int           `bson:"country_id" json:"country_id"`
	RomGroupId  bson.ObjectId `bson:"rom_group_id,omitempty" json:"rom_group_id" sql:"-"`

	// Could be used for specific rom art otherwise fall back to rom-group
	BoxArts []BoxArt `bson:"box_arts" json:"box_arts"`
}

type BoxArt struct {
	Height string `json:"height"`
	Width  string `json:"width"`
	Src    string `json:"src"`
	Side   string `json:"side"`
}

func NewRom(goodName string) *Rom {
	rom := &Rom{}
	rom.GoodName = goodName
	rom.System = 0

	return rom
}

func (self *Rom) String() string {
	return fmt.Sprintf(`Name: %s
Good: %s
Md5: %s`, self.Name, self.GoodName, self.Md5)
}
