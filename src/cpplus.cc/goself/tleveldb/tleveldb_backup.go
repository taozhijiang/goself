package tleveldb


import (
	"os"
	"bufio"
    "log"
    b64 "encoding/base64"

    "strings"
    "errors"

    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/opt"
)

func BrowseLevelDB(path string) (uint64, error) {

	o := &opt.Options{
		ReadOnly: true,
	}
	db, err := leveldb.OpenFile(path, o)
	if err != nil {
		log.Printf("Open(ReadOnly) %s failed.", path)
		return 0, err
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	var count uint64 = 0
	for iter.Next() {  // may break early when error
		key := iter.Key()
		value := iter.Value()

		log.Printf("%d => %s: %d", count, string(key), len(value))
		count ++
	}

	iter.Release()
	err = iter.Error()
	if err != nil {
		log.Print("Iter error:", err)
		return 0, err
	}

	return count, nil
}

func BackupLevelDB(path string, backtar string) (uint64, error) {

	// open backfile
	f, err := os.OpenFile(backtar, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644) // default O_RDONLY
	if err != nil {
		log.Print("Error when open: ", backtar, err)
		return 0, err
	}
	defer f.Close()

	o := &opt.Options{
		ReadOnly: true,
		ErrorIfMissing: true,
	}
	db, err := leveldb.OpenFile(path, o)
	if err != nil {
		log.Printf("Open LevelDB(Read) %s failed.", path)
		return 0, err
	}
	defer db.Close()

	var count uint64 = 0
	iter := db.NewIterator(nil, nil)

	for iter.Next() {  // may break early when error
		key := iter.Key()
		value := iter.Value()
		log.Print("Back:", string(key))

		sEnckey   := b64.StdEncoding.EncodeToString([]byte(key))
		sEncvalue := b64.StdEncoding.EncodeToString([]byte(value))

		content := sEnckey + ":" + sEncvalue + "\n"
		f.Write([]byte(content))
		count ++
	}

	iter.Release()
	err = iter.Error()
	if err != nil {
		log.Print("Iter error.", err)
		return 0, err
	}

	return count, nil
}



func RestoreLevelDB(path string, backtar string) (uint64, error) {

	// open backupfile
	f, err := os.Open(backtar) // default O_RDONLY
	if err != nil {
		log.Print("Error when open: ", backtar)
		return 0, err
	}
	defer f.Close()

	o := &opt.Options{
		ReadOnly: false,
		ErrorIfExist: true,
	}
	db, err := leveldb.OpenFile(path, o)
	if err != nil {
		log.Printf("Create LevelDB %s failed.", path)
		return 0, err
	}
	defer db.Close()

	// batch restore
	batch := new(leveldb.Batch)
	var count uint64 = 0

	input := bufio.NewScanner(f)
	buf := make([]byte, 0, 1024*1024) // Initial 1M
	input.Buffer(buf, 100*1024*1024)  // Maxium 100M
	for input.Scan() {
		backstr := input.Text()
		backitem := strings.Split(backstr, ":")
		if len(backitem) != 2 {
            log.Print("Error process:", backstr)
			return 0, errors.New("Format error")
		}

		sDeckey, err := b64.StdEncoding.DecodeString(backitem[0])
		if err != nil {
			return 0, err
		}
		sDecval, err := b64.StdEncoding.DecodeString(backitem[1])
		if err != nil {
			return 0, err
		}

		log.Print("Restore:", string(sDeckey))
		batch.Put([]byte(sDeckey), []byte(sDecval))
		count ++

        // Commit batch write ops every 10 items
		if count >= 10 {
			err = db.Write(batch, nil)
			if err != nil {
				log.Print("Write to LevelDB error.", err)
				return 0, err
			}

			batch = new(leveldb.Batch)
		}

	}

	if err = input.Err(); err != nil {
        log.Printf("Invalid input: %s", err)
        return 0, err
	}

	err = db.Write(batch, nil)
	if err != nil {
	    log.Print("Write to LevelDB error.", err)
		return 0, err
	}

	return count, nil
}
