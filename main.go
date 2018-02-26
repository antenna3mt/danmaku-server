// Copyright 2018 Yi Jin. All rights reserved.
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"github.com/antenna3mt/rpc"
	"github.com/antenna3mt/rpc/json"
	"net/http"
	"fmt"
	"github.com/rs/cors"
)

func main() {
	engine := NewEngine()
	engine.NewActivityFull(engine.AdminToken, "Test Activity", "cc123456", "rr123456", "dd123456")
	fmt.Println(engine.AdminToken)
	server, err := rpc.NewServer(new(Context))
	if err != nil {
		log.Fatal(err)
	}
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterService(&DanmakuService{
		E: engine,
	}, "")
	http.Handle("/", cors.Default().Handler(server))
	if err := http.ListenAndServe(":8881", nil); err != nil {
		log.Fatal(err)
	}
}
