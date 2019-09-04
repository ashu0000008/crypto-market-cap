package server

import (
	"container/list"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func StartWSServer(chanCollector chan string, chanData chan string) {

	initRelation()
	go processData(chanData)

	http.ListenAndServe(":9000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}

		go func() {
			defer func() {
				conn.Close()
				RemoveConn(conn)
			}()

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					// handle error
				}

				if op == ws.OpContinuation {
					time.Sleep(time.Duration(1) * time.Second)
				}

				proccessRequest(string(msg), chanCollector, conn)

				//echo := "server get:" + string(msg)
				//err = wsutil.WriteServerMessage(conn, op, []byte(echo))
				//if err != nil {
				//	// handle error
				//}
			}
		}()
	}))
}

func proccessRequest(request string, chanCollector chan string, conn net.Conn) {
	tmp := strings.Split(request, ":")
	if len(tmp) != 2 {
		return
	}

	symbol := tmp[1]
	chanCollector <- symbol

	//将chan与symbol的关系通知给manager
	if strings.EqualFold("add", tmp[0]) {
		AddRelation(conn, symbol)
	} else if strings.EqualFold("remove", tmp[0]) {
		RemoveRelation(conn, symbol)
	}
}

func processData(chanData chan string) {
	var input string
	for {
		input = <-chanData
		processDataIndeed(input)
	}
}

type relation struct {
	conn   net.Conn
	symbol string
}

type relations struct {
	data *list.List
	lock *sync.Mutex
}

var _relations relations

func initRelation() {
	_relations.data = list.New()
	_relations.lock = &sync.Mutex{}
}

func AddRelation(conn1 net.Conn, symbol1 string) {
	fmt.Println("AddRelation", conn1, symbol1)

	_relations.lock.Lock()

	data := relation{conn: conn1, symbol: symbol1}
	_relations.data.PushBack(data)

	_relations.lock.Unlock()
}

func RemoveRelation(conn net.Conn, symbol string) {
	fmt.Println("RemoveRelation", conn, symbol)

	_relations.lock.Lock()

	for e := _relations.data.Front(); e != nil; e = e.Next() {
		if e.Value.(*relation).conn == conn && e.Value.(*relation).symbol == symbol {
			_relations.data.Remove(e)
			break
		}
	}

	_relations.lock.Unlock()
}

func RemoveConn(conn net.Conn) {
	fmt.Println("RemoveRelation", conn)

	_relations.lock.Lock()

	for e := _relations.data.Front(); e != nil; e = e.Next() {
		if e.Value.(*relation).conn == conn {
			_relations.data.Remove(e)
			break
		}
	}

	_relations.lock.Unlock()
}

func processDataIndeed(data string) {
	fmt.Println("processDataIndeed", data)

	tmp := strings.Split(data, "-")
	if len(tmp) != 2 {
		return
	}
	symbol := tmp[0]

	_relations.lock.Lock()

	for e := _relations.data.Front(); e != nil; e = e.Next() {
		if strings.Contains(e.Value.(relation).symbol, symbol) {
			err := wsutil.WriteServerMessage(e.Value.(relation).conn, ws.OpContinuation, []byte(data))
			if err != nil {
				// handle error
			}
		}
	}

	_relations.lock.Unlock()
}
