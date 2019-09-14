package server

import (
	"container/list"
	"fmt"
	"github.com/ashu0000008/crypto-market-cap/ws/config"
	"github.com/ashu0000008/crypto-market-cap/ws/redisops"
	"github.com/go-redis/redis"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func StartWSServer() {

	redisClient := redisops.RedisConnect()
	dataPubsub := redisClient.Subscribe(config.REDIS_QUOTE_DATA_NAME)

	initRelation()
	go processData(dataPubsub)

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
					fmt.Println(err.Error())

					//链接坏了，就退出
					if processNetError(conn, err) {
						break
					}
				}

				if op == ws.OpClose {
					fmt.Println("conn", conn, "close")
					break
				}

				if op == ws.OpContinuation {
					time.Sleep(time.Duration(1) * time.Second)
				}

				fmt.Println("conn receive--------:", string(msg))
				proccessRequest(string(msg), redisClient, conn)

				//echo := "server get:" + string(msg)
				//err = wsutil.WriteServerMessage(conn, op, []byte(echo))
				//if err != nil {
				//	// handle error
				//}
			}

			fmt.Println("-------quit", conn)
		}()
	}))
}

func processNetError(conn net.Conn, err error) bool {

	//netErr, ok := err.(net.Error)
	//if !ok {
	//	return false
	//}
	//
	//if netErr.Timeout() {
	//	log.Println("timeout")
	//	return false
	//}

	opErr, ok := err.(*net.OpError)
	if !ok {
		return false
	}

	switch t := opErr.Err.(type) {
	case *net.DNSError:
		log.Printf("net.DNSError:%+v", t)
		return true
	case *os.SyscallError:
		log.Printf("os.SyscallError:%+v", t)
		if errno, ok := t.Err.(syscall.Errno); ok {
			switch errno {
			case syscall.ECONNREFUSED:
				log.Println("connect refused")
				return true
			case syscall.ECONNRESET:
				log.Println("connect reset")
				return true
			case syscall.ETIMEDOUT:
				log.Println("timeout")
				return false
			}
		}
	}

	return false
}

func proccessRequest(request string, redisClient *redis.Client, conn net.Conn) {
	tmp := strings.Split(request, ":")
	if len(tmp) != 2 {
		return
	}

	symbol := tmp[1]
	redisClient.Publish(config.REDIS_QUOTE_MANAGER_NAME, symbol)

	//将chan与symbol的关系通知给manager
	if strings.EqualFold("add", tmp[0]) {
		AddRelation(conn, symbol)
	} else if strings.EqualFold("remove", tmp[0]) {
		RemoveRelation(conn, symbol)
	}
}

func processData(dataPubsub *redis.PubSub) {
	for {
		input, _ := dataPubsub.ReceiveMessage()
		processDataIndeed(input.Payload)
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

	data := &relation{conn: conn1, symbol: symbol1}
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

	var n *list.Element
	for e := _relations.data.Front(); e != nil; e = n {
		n = e.Next()
		if reflect.DeepEqual(e.Value.(*relation).conn, conn) {
			_relations.data.Remove(e)
		}
	}

	_relations.lock.Unlock()
}

func processDataIndeed(data string) {
	//fmt.Println("processDataIndeed", data)

	tmp := strings.Split(data, "-")
	if len(tmp) != 2 {
		return
	}
	symbol := tmp[0]

	_relations.lock.Lock()

	for e := _relations.data.Front(); e != nil; e = e.Next() {
		target := e.Value.(*relation).symbol
		if strings.Contains(symbol, target) {
			fmt.Println("send data:", data)
			err := wsutil.WriteServerMessage(e.Value.(*relation).conn, ws.OpText, []byte(data))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	_relations.lock.Unlock()
}
