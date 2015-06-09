package common

import (
	_ "fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Database struct {
	session       *mgo.Session
	RomCollection *mgo.Collection
}

func NewDatabase() *Database {
	importer := new(Database)
	session, err := mgo.Dial("mongodb://localhost")
	db := "goodnes"

	if err != nil {
		panic(err)
	}
	importer.session = session
	importer.RomCollection = session.DB(db).C("roms")

	return importer
}
func (self *Database) PersistRom(rom *Rom) error {
	dRom := new(Rom)
	err := self.RomCollection.Find(bson.M{"md5": rom.Md5}).One(&dRom)
	if err == mgo.ErrNotFound {
		self.RomCollection.Insert(rom)
	} else {
		self.RomCollection.Update(bson.M{"md5": rom.Md5}, rom)
		if err != nil {
			log.Fatal("Error querying:", err)
		}
	}
	return nil
}

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

type Rom struct {
	Name     string `json:"name"`
	GoodName string `bson:"good_name" json:"sanitized_name"`
	Md5      string `json:"md5"`
	System   int8   `json:"system"`

	Hack        bool   `json:"hack"`
	HackName    string `bson:"hack_name" json:"hack_name"`
	Trainer     bool   `json:"trainer"`
	TrainerName string `bson:"trainer_name" json:"trainer_name"`
	BadDump     bool   `bson:"bad_dump" json:"bad_dump"`
	Overdump    bool   `json:"overdump"`
	Verified    bool   `json:"verified"`
	Alternative bool   `json:"alternative"`
	Pirated     bool   `json:"pirated"`
	Fixed       bool   `json:"fixed"`
	CountryId   int    `bson:"country_id" json:"country_id"`

	GamesDbId      int  `bson:"games_db_id" json:"games_db_id"`
	GamesDbScanned bool `bson:"initial_games_db_scan" json:"games_db_scanned"`

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
