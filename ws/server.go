package main

import (
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			// handle error
		}
		go func() {
			defer conn.Close()

			for {
				msg, op, err := wsutil.ReadClientData(conn)
				if err != nil {
					// handle error
				}

				if op == ws.OpContinuation {
					time.Sleep(time.Duration(1) * time.Second)
				}

				echo := "server get:" + string(msg)
				err = wsutil.WriteServerMessage(conn, op, []byte(echo))
				if err != nil {
					// handle error
				}
			}
		}()
	}))
}
