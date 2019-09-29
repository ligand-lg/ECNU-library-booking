package main

import (
	"log"
	"time"
)

/*
为了防止本程序在后台无限循环，把学校服务器搞蹦，设置一个计时器，计时器一到，程序自杀。
*/

const (
	LifeTime = 10 * time.Second
)

var _suicideBoom = time.After(1 * time.Second)
var _suicideInit = false

func suicide() {
	if !_suicideInit {
		_suicideBoom = time.After(LifeTime)
		_suicideInit = true
		return
	} else {
		select {
		case <-_suicideBoom:
			log.Fatal("爆炸启动，BOOM...")
		default:
			return
		}
	}
}
