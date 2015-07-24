package gen

import "math/rand"

var _chunkMigrationGen *ChunkMigrationGen

type ChunkMigrationGen struct {
	C chan [3]int64 // from -> to
}

func (r *ChunkMigrationGen) Init() {
	r.C = make(chan [3]int64)

	go func() {
		for {
			<-r.C
			r.C <- [3]int64{rand.Int63(), rand.Int63n(3), rand.Int63n(3)}
		}
	}()
}

func ChunkMigration() *ChunkMigrationGen {
	return _chunkMigrationGen
}

func init() {
	_chunkMigrationGen = &ChunkMigrationGen{}
	_chunkMigrationGen.Init()
}
