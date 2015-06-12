package common

import (
	_ "fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
)

type Database struct {
	session            *mgo.Session
	RomCollection      *mgo.Collection
	RomGroupCollection *mgo.Collection
}

func NewDatabase() *Database {
	importer := new(Database)
	logout := log.New(os.Stdout, "[mGo]: ", log.Lshortfile)
	mgo.SetLogger(logout)
	mgo.SetDebug(false)
	session, err := mgo.Dial("mongodb://localhost")
	db := "goodnes"

	if err != nil {
		panic(err)
	}
	importer.session = session
	importer.RomCollection = session.DB(db).C("roms")
	importer.RomGroupCollection = session.DB(db).C("rom_groups")

	return importer
}
func (self *Database) FindOrCreateGroup(name string) *RomGroup {
	var romGroup RomGroup
	err := self.RomGroupCollection.Find(bson.M{"name": name}).One(&romGroup)
	if err == mgo.ErrNotFound {
		log.Println("[db] Romgroup not found, creating.")
		romGroup.Name = name
		self.RomGroupCollection.Insert(romGroup)
	} else {
		log.Println("[db] Romgroup found", romGroup.Id)
	}
	return &romGroup
}
func (self *Database) PersistRomGroup(romGroup *RomGroup) bool {
	err := self.RomGroupCollection.Update(bson.M{"_id": romGroup.Id}, romGroup)
	if err != nil {
		log.Println("[db] Error saving romGroup")
		return false
	} else {
		return true
	}
}

func (self *Database) FindOrInitializeRom(md5 string) *Rom {
	var rom Rom
	err := self.RomCollection.Find(bson.M{"md5": md5}).One(&rom)
	if err == mgo.ErrNotFound {
		log.Println("[db] Rom not found yet, initializing.")
		rom.Md5 = md5
	}
	return &rom
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
