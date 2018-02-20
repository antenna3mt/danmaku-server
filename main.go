// Copyright 2018 Yi Jin. All rights reserved.
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"github.com/antenna3mt/rpc"
	"github.com/antenna3mt/rpc/json"
	"net/http"
	"fmt"
)

func main() {
	engine := NewEngine()
	fmt.Println(engine.AdminToken)
	server, err := rpc.NewServer(new(Context))
	if err != nil {
		log.Fatal(err)
	}
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterService(&DanmakuService{
		E: engine,
	}, "")
	http.Handle("/", server)
	if err := http.ListenAndServe(":8881", nil); err != nil {
		log.Fatal(err)
	}
}
