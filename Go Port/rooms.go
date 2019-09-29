package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

/**
中北 + 闵行校区小黑屋配置。
*/

type __Rooms struct {
	Rooms []Room `json:"data"`
}

type Room struct {
	DevId   string `json:"devId"`
	KindId  string `json:"kindId"`
	LabId   string `json:"labId"`
	DevName string `json:"devName"`
	RoomNo  string `json:"roomNo"`
}

func initRoom() []Room {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	zbContent, err := ioutil.ReadFile(filepath.Join(dir, "zb.json"))
	//zbContent, err := ioutil.ReadFile("zb.json")
	if err != nil {
		log.Fatal(err)
	}
	mhContent, err := ioutil.ReadFile(filepath.Join(dir, "mh.json"))
	//mhContent, err := ioutil.ReadFile("mh.json")
	if err != nil {
		log.Fatal(err)
	}
	var zbJson __Rooms
	var mhJson __Rooms
	err = json.Unmarshal(zbContent, &zbJson)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(mhContent, &mhJson)
	if err != nil {
		log.Fatal(err)
	}

	res := append(zbJson.Rooms, mhJson.Rooms...)
	return res
}

func GetRoom(roomNo string) (*Room, error) {
	__rooms := initRoom()
	for _, v := range __rooms {
		if v.RoomNo == roomNo {
			return &v, nil
		}
	}
	return nil, errors.New("找不到房间号为:" + roomNo + "的房间。")
}
