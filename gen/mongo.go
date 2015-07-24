package gen

import (
	"encoding/json"
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*
	"opcounters" : {
		"insert" : 1,
		"query" : 17,
		"update" : 0,
		"delete" : 0,
		"getmore" : 914,
		"command" : 899
	},
*/
type MongoServerStatus struct {
	UptimeMillis map[string]int64 `json:"uptimeMillis"`
	Opcounters   map[string]int64 `json:"opcounters"`
}

func getMongoServerStatus(s string) *MongoServerStatus {
	r := &MongoServerStatus{}

	err := json.Unmarshal([]byte(s), &r)
	if err != nil {
		log.Fatalf("Failed to parse serverStatus with error %s\n\n", fmt.Sprint(err), s)
	}

	return r
}

func getMongoServerStatus_old(s *mgo.Session) *MongoServerStatus {
	r := &MongoServerStatus{}
	err := s.Run(bson.D{{"serverStatus", 1}}, r)

	if err != nil {
		log.Fatalf("could not run serverStatus, got error: ", err)
	}

	return r
}
