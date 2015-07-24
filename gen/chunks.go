package gen

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _mongosMonitor *mongosMonitor

type MongosStats struct {
	LastLine string

	ChunkDistribution [3]int
}

type mongosMonitor struct {
	S MongosStats
}

func (r *mongosMonitor) Init() {

	// start monitor
	r.Run()
}

func (r *mongosMonitor) GetStats() MongosStats {
	// FIXME do I need lock here?
	return r.S
}

func (s *mongosMonitor) Run() {

	// check mongod/s information

	dialInfo := mgo.DialInfo{
		FailFast: true,
		Addrs:    []string{Getenv("ms") + ":58989"},
	}

	m, err := mgo.DialWithInfo(&dialInfo)

	if err != nil {
		log.Fatalf("Failed to dial server ms/%s with error %s\n", dialInfo.Addrs[0], fmt.Sprint(err))
	}

	go func() {
		for {
			for i := 0; i < 3; i++ {
				s.S.ChunkDistribution[i], err = m.DB("config").C("chunks").Find(bson.M{"shard": fmt.Sprintf("rs%1d", i)}).Count()

				if err != nil {
					log.Fatalf("Cannot count chunk for shard r%1d", i)
				}
			}

			time.Sleep(1 * time.Second)
		}
	}()
}

func MongosMonitor() *mongosMonitor {
	return _mongosMonitor
}

func init() {
	_mongosMonitor = &mongosMonitor{}
	_mongosMonitor.Init()
}
