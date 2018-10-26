package tleveldb

import (
	"log"
	"errors"
    "testing"

    "github.com/syndtr/goleveldb/leveldb"
)

var dbpath string = "./testdb"
var dbpath2 string = "./testdb2"
var dbback string = "./testdb.bin"

func writeDB(path string) error {

	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Printf("Open(Read) %s failed.", path)
		return err
	}
	defer db.Close()

	if err = db.Put([]byte("key111"), []byte("val111"), nil); err != nil {
		return err
	}

	if err = db.Put([]byte("key222"), []byte("val222"), nil); err != nil {
		return err
	}

	return nil
}

func checkDB(path string) error {

	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Printf("Open(Read) %s failed.", path)
		return err
	}
	defer db.Close()

	var value []byte

	value, err = db.Get([]byte("key111"), nil)
	if string(value) != "val111" || err != nil {
		log.Print("==>", string(value))
		return errors.New("CheckError")
	}

	value, err = db.Get([]byte("key222"), nil)
	if string(value) != "val222" || err != nil {
		log.Print("==>", string(value))
		return errors.New("CheckError")
	}

	return nil
}


func TestLevelDB(t *testing.T) {
	
	var err error

	if err = writeDB(dbpath); err != nil {
		t.Fatal(err)
	}

	cnt, err := BackupLevelDB(dbpath, dbback);
	if cnt != 2 || err != nil {
		t.Fatal("error when BackupLevelDB:", cnt, err)
	}

	err = checkDB(dbpath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = BrowseLevelDB(dbpath)
	if err != nil {
		t.Fatal(err)
	}

	cnt, err = RestoreLevelDB(dbpath2, dbback)
	if cnt != 2 || err != nil {
		t.Fatal("error when RestoreLevelDB:", cnt, err)
	}

	_, err = BrowseLevelDB(dbpath2)
	if err != nil {
		t.Fatal(err)
	}
}