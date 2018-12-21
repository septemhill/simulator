package simulator

import "github.com/syndtr/goleveldb/leveldb"

func OpenDatabase(path string) {
	leveldb.OpenFile(path, nil)
}
