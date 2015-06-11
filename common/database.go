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
