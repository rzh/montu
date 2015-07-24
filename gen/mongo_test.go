package gen

import (
	"fmt"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestMongo_serverStats(t *testing.T) {

	dialInfo := mgo.DialInfo{
		FailFast: true,
		Addrs:    []string{"54.148.50.239:58989"},
	}

	s, err := mgo.DialWithInfo(&dialInfo)

	if err != nil {
		t.Errorf("Failed to dial to server")
	}

	r := &MongoServerStatus{}
	err = s.Run(bson.D{{"serverStatus", 1}}, r)

	if err != nil {
		t.Error("could not run serverStatus, got error: ", err)
	}

	fmt.Println(r)
}

func TestMongo_serverStats_old(t *testing.T) {

}

func TestServerMinotor(t *testing.T) {
	m := MonitorServer("p3")

	fmt.Printf("%+v", m.GetStats().mongo_stats)

	time.Sleep(2 * time.Second)

	fmt.Printf("%+v", m.GetStats().mongo_stats)
}
