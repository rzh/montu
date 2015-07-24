package gen

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ServerStat struct {
	LastLine       string
	P_CPU          int
	P_MEM          int
	P_DISK_dbs     int
	P_DISK_journal int
	TS             int

	OPS         int
	OPS_insert  int
	OPS_update  int
	OPS_query   int
	OPS_getmore int
	OPS_delete  int
	OPS_command int

	// internal
	mongo_stats *MongoServerStatus
}

var _serverMonitors map[string]*ServerMonitor
var _mutex *sync.Mutex = &sync.Mutex{}

type ServerMonitor struct {
	C      chan ServerStat // from -> to
	S      ServerStat
	server string // this shall be the env var for the server, ms, p1, etc.
}

func Getenv(e string) (s string) {
	s = os.Getenv(e)

	if s == "" {
		log.Panicf("Failed to read address for server %s\n", e)
	}

	return
}

func (r *ServerMonitor) Init(addr string) {
	r.C = make(chan ServerStat)
	r.server = addr

	// start monitor
	r.Run()

	go func() {
		for {
			<-r.C
			r.C <- r.S
		}
	}()
}

func (r *ServerMonitor) GetStats() ServerStat {
	// FIXME do I need lock here?
	return r.S
}

func cleanString(l string) string {
	var t string
	var last_space bool
	last_space = false

	l = strings.Replace(l, "|", " ", -1)
	l = strings.Replace(l, ":", " ", -1)

	for i := 0; i < len(l); i++ {

		if l[i] == ' ' && last_space {
			// ignore duplicate space
		} else if l[i] == ' ' {
			last_space = true
			t = t + " "
		} else {
			last_space = false
			t = t + l[i:i+1]
		}
	}

	return strings.TrimSpace(t)
}

func getMemory(s string) int64 {
	// 100M -> 100
	// 100G -> 100 * 1024

	if s[len(s)-1] == 'M' {
		a, err := strconv.ParseFloat(s[0:len(s)-1], 64)

		if err != nil {
			log.Fatalf("Cannot conver memory %s\n", s)
		}
		return int64(a)
	} else if s[len(s)-1] == 'G' {
		a, err := strconv.ParseFloat(s[0:len(s)-1], 64)

		if err != nil {
			log.Fatalf("Cannot conver memory %s\n", s)
		}
		return int64(a * 1024)
	}
	return 0
}

func (s *ServerMonitor) Parse(l string) {
	ss := strings.Split(cleanString(l), " ")
	_ = ss[0]
	// sample line
	//  --epoch--- ----total-cpu-usage---- ------memory-usage----- xvda -net/total-
	//    epoch   |usr sys idl wai hiq siq| used  buff  cach  free|util| recv  send
	//  1437509148|  0   0 100   0   0   0|2872M  201M 5466M 50.6G|0.01|   0     0
	var err error

	// MEM
	used := getMemory(ss[7])
	buff := getMemory(ss[8])
	cach := getMemory(ss[9])
	free := getMemory(ss[10])

	s.S.P_MEM = 100 - int((100*free)/(free+used+buff+cach))

	// CPU
	s.S.P_CPU, err = strconv.Atoi(strings.Trim(ss[3], " \t"))
	if err != nil {
		log.Fatalln("Failed to parse CPU ", ss, ss[3], " from dstat --> ", err, " | ", l)
	}
	s.S.P_CPU = 100 - s.S.P_CPU // adjust to show real CPU usage

	// DISK
	if len(ss) == 16 {
		// mongod

		if f, err := strconv.ParseFloat(ss[14], 64); err != nil {
			log.Fatalf("Failed to parse float %s\n", ss[14])
		} else {
			s.S.P_DISK_dbs = int(f)
		}

		if f, err := strconv.ParseFloat(ss[15], 64); err != nil {
			log.Fatalf("Failed to parse float %s\n", ss[15])
		} else {
			s.S.P_DISK_journal = int(f)
		}
	} else {
		if f, err := strconv.ParseFloat(ss[13], 64); err != nil {
			log.Fatalf("Failed to parse float %s\n", ss[13])
		} else {
			s.S.P_DISK_dbs = int(f)
			s.S.P_DISK_journal = 0
		}
	}
}

type CMD string

func (s *ServerMonitor) Run() {

	go func() {
		Cmd_exec := exec.Command(
			"/usr/bin/ssh",
			"-o", "StrictHostKeyChecking=no",
			"-i", "/Users/rui/bin/rui-aws-cap.pem",
			"ec2-user@"+Getenv(s.server),
			"dstat -T -c -m -n --disk-util")

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)

		go func() {
			// Block until a signal is received.
			<-c
			Cmd_exec.Process.Kill()
		}()

		out, e := Cmd_exec.StdoutPipe()
		if e != nil {
			log.Fatalln("Cannot get StdoutPipe with error:", e)
		}
		r := bufio.NewReader(out)

		if err := Cmd_exec.Start(); err != nil {
			log.Fatal("failed with -> ", err)
		}

		var err error
		var l string

		// skip the first two line
		l, err = r.ReadString('\n')
		l, err = r.ReadString('\n')

		for err == nil {
			l, err = r.ReadString('\n')
			l = strings.TrimRight(l, "\n\t ")
			l = cleanString(l)

			if l != "" {
				s.S.LastLine = l

				// parse
				s.Parse(l)
			}
		}

		Cmd_exec.Wait()
	}()

	// check mongod/s information

	/*
		dialInfo := mgo.DialInfo{
			FailFast: true,
			Addrs:    []string{Getenv(s.server) + ":58989"},
		}

		m, err := mgo.DialWithInfo(&dialInfo)

		if err != nil {
			log.Fatalf("Failed to dial server %s/%s with error %s\n", s.server, dialInfo.Addrs[0], fmt.Sprint(err))
		}

	*/

	go func() {
		Cmd_exec := exec.Command(
			"/usr/bin/ssh",
			"-o", "StrictHostKeyChecking=no",
			"-i", "/Users/rui/bin/rui-aws-cap.pem",
			"ec2-user@"+Getenv(s.server),
			"~/bin/mongo 127.0.0.1:58989 --eval \"while(true) {print(JSON.stringify(db.serverStatus())); sleep(1000)}\"")

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)

		go func() {
			// Block until a signal is received.
			<-c
			Cmd_exec.Process.Kill()
		}()

		out, e := Cmd_exec.StdoutPipe()
		if e != nil {
			log.Fatalln("Cannot get StdoutPipe with error:", e)
		}
		r := bufio.NewReader(out)

		if err := Cmd_exec.Start(); err != nil {
			log.Fatal("failed with -> ", err)
		}
		var err error

		// skip the first two line
		_, err = r.ReadString('\n')
		_, err = r.ReadString('\n')

		for err == nil {
			m, err := r.ReadString('\n')

			if err != nil {
				log.Fatalf("Failed to read server status for %s\n", s.server)
			}
			mongo_stats := getMongoServerStatus(m)

			var totalOps = func(o *MongoServerStatus) (ops int64) {
				for _, value := range o.Opcounters {
					ops = ops + value
				}
				return
			}

			if s.S.mongo_stats != nil {
				s.S.OPS = int((totalOps(mongo_stats) - totalOps(s.S.mongo_stats)) * 1000 / (mongo_stats.UptimeMillis["floatApprox"] - s.S.mongo_stats.UptimeMillis["floatApprox"]))

				// insert
				ops_type := "insert"
				t := int((mongo_stats.Opcounters[ops_type] - s.S.mongo_stats.Opcounters[ops_type]) * 1000 / (mongo_stats.UptimeMillis["floatApprox"] - s.S.mongo_stats.UptimeMillis["floatApprox"]))
				s.S.OPS_insert = t
				// update
				ops_type = "update"
				t = int((mongo_stats.Opcounters[ops_type] - s.S.mongo_stats.Opcounters[ops_type]) * 1000 / (mongo_stats.UptimeMillis["floatApprox"] - s.S.mongo_stats.UptimeMillis["floatApprox"]))
				s.S.OPS_update = t
				// query
				ops_type = "query"
				t = int((mongo_stats.Opcounters[ops_type] - s.S.mongo_stats.Opcounters[ops_type]) * 1000 / (mongo_stats.UptimeMillis["floatApprox"] - s.S.mongo_stats.UptimeMillis["floatApprox"]))
				s.S.OPS_query = t
				// command
				ops_type = "command"
				t = int((mongo_stats.Opcounters[ops_type] - s.S.mongo_stats.Opcounters[ops_type]) * 1000 / (mongo_stats.UptimeMillis["floatApprox"] - s.S.mongo_stats.UptimeMillis["floatApprox"]))
				s.S.OPS_command = t
				// delete
				ops_type = "delete"
				t = int((mongo_stats.Opcounters[ops_type] - s.S.mongo_stats.Opcounters[ops_type]) * 1000 / (mongo_stats.UptimeMillis["floatApprox"] - s.S.mongo_stats.UptimeMillis["floatApprox"]))
				s.S.OPS_delete = t
				// getmore
				ops_type = "getmore"
				t = int((mongo_stats.Opcounters[ops_type] - s.S.mongo_stats.Opcounters[ops_type]) * 1000 / (mongo_stats.UptimeMillis["floatApprox"] - s.S.mongo_stats.UptimeMillis["floatApprox"]))
				s.S.OPS_getmore = t
			}
			s.S.mongo_stats = mongo_stats

			time.Sleep(1 * time.Second)
		}
	}()
}

func MonitorServer(server_url string) *ServerMonitor {
	_mutex.Lock()

	if _, ok := _serverMonitors[server_url]; ok {
		_mutex.Unlock()
		return _serverMonitors[server_url]
	}

	m := ServerMonitor{}
	m.Init(server_url)
	_serverMonitors[server_url] = &m

	_mutex.Unlock()

	return _serverMonitors[server_url]
}

func init() {
	_serverMonitors = make(map[string]*ServerMonitor)
}
