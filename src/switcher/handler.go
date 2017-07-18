package main

import (
	"comm"
	"comm/config"
	"comm/sched/loop"
	"encoding/json"
	"fmt"
	"net/http"
	"switcher/app"
	"switcher/app/dbmgr"
)

func HandlerGateReg(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	token, ok1 := r.Form["token"]
	name, ok2 := r.Form["name"]
	ip, ok3 := r.Form["ip"]
	port, ok4 := r.Form["port"]

	if !ok1 || !ok2 || !ok3 || !ok4 {
		fmt.Fprint(w, "need param token, name, ip, port")
		return
	}

	if token[0] != config.Switcher.Token {
		fmt.Fprint(w, "token error")
		return
	}

	loop.Push(func() {
		app.RegGate(name[0], ip[0], port[0])
	})

	fmt.Fprint(w, "ok")
}

func HandlerGameList(w http.ResponseWriter, r *http.Request) {
	res := []string{}

	for k, _ := range config.Games {
		res = append(res, k)
	}

	bytes, _ := json.Marshal(res)
	fmt.Fprint(w, string(bytes))
}

func HandlerGateList(w http.ResponseWriter, r *http.Request) {
	app.PrintGates()

	jg, err := json.Marshal(app.Gates)
	if err != nil {
		log.Error("HandlerGateList Marshal failed")
		return
	}

	fmt.Fprint(w, string(jg))
}

func HandlerUserSvr(w http.ResponseWriter, r *http.Request) {
	//client need send param both 'channel_uid' and 'svr0' to switcher
	r.ParseForm()
	if r.Method == "POST" {
		channel := r.PostFormValue("channel")
		uid := r.PostFormValue("channel_uid")
		svr0 := r.PostFormValue("svr0")
		log.Info("uid, svr0:", uid, svr0)
		if len(channel) == 0 || len(uid) == 0 || len(svr0) == 0 {
			fmt.Fprint(w, "falied")
		}

		svr := dbmgr.CenterGetUserSvr(comm.Atoi32(channel), uid, svr0)
		fmt.Fprint(w, fmt.Sprintf("{\"svr\":\"%s\"}", svr))
	} else {
		fmt.Fprint(w, "falied")
	}
}
