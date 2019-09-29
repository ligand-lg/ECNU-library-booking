package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	MaxDuration = 4 * time.Hour
	MinDuration = 30 * time.Minute
	MinInterval = 10 * time.Minute
)

// 对应一条预定
type Booking struct {
	room Room
	// 开发时测试用
	delayDay int
	// 预定开始时间
	startTime string
	// 预定结束时间
	endTime string
}

// 总配置
type Config struct {
	sid        string
	pwd        string
	vip        bool
	allBooking []Booking
}

// vip 优先预定
type _config struct {
	Sid      string      `json:"sid"`
	Password string      `json:"password"`
	RoomNo   string      `json:"roomNo"`
	Vip      bool        `json:"vip"`
	Duration [][2]string `json:"duration"`
}

func getConfFilePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(dir, "conf.json")
	//return "conf.json"
}

// ref: https://yourbasic.org/golang/format-parse-string-time-date-example/
func TimeInterval(startStr, endStr string) time.Duration {
	layout := "15:04"
	start, err1 := time.Parse(layout, startStr)
	end, err2 := time.Parse(layout, endStr)
	if err1 != nil {
		log.Println(err1)
	}
	if err2 != nil {
		log.Println(err2)
	}
	return end.Sub(start)
}

func checkConf(conf Config) (pass bool) {
	pass = true
	if len(conf.sid) == 0 {
		pass = false
		log.Println("学号不能为空")
	}
	if len(conf.pwd) == 0 {
		pass = false
		log.Println("密码不能为空")
	}
	if len(conf.allBooking) == 0 {
		pass = false
		log.Println("预定不能为空")
	}
	for _, b := range conf.allBooking {
		if len(b.startTime) != 5 {
			pass = false
			log.Println("错误的时间格式：", b.startTime)
		}
		if len(b.endTime) != 5 {
			pass = false
			log.Println("错误的时间格式：", b.endTime)
		}
	}
	// 预定校验
	sort.SliceStable(conf.allBooking, func(i, j int) bool {
		return conf.allBooking[i].startTime < conf.allBooking[j].endTime
	})
	preEnd := "07:50"
	for _, d := range conf.allBooking {
		if TimeInterval(preEnd, d.startTime) < MinInterval {
			pass = false
			log.Println("最短预约间隔不能小于：", MinInterval)
			return
		}
		ddd := TimeInterval(d.startTime, d.endTime)
		if ddd < MinDuration {
			pass = false
			log.Println("预定持续时间不能小于：", MinDuration)
			return
		}
		if ddd > MaxDuration {
			pass = false
			log.Println("预定持续时间不能大于：", MaxDuration)
			return
		}
		preEnd = d.endTime
	}
	return
}

func GetConf() (conf Config) {
	byteValue, err := ioutil.ReadFile(getConfFilePath())
	if err != nil {
		log.Fatal(err)
	}
	var jsonObj _config
	err = json.Unmarshal(byteValue, &jsonObj)
	if err != nil {
		log.Fatalln(err)
	}
	room, err := GetRoom(jsonObj.RoomNo)
	if err != nil {
		log.Fatalln(err)
	}
	allBooking := make([]Booking, len(jsonObj.Duration))

	delayDay := 2
	if DEBUG {
		delayDay = 1
	}
	for i, d := range jsonObj.Duration {
		allBooking[i] = Booking{
			room: *room,
			// 默认定为后天
			delayDay:  delayDay,
			startTime: d[0],
			endTime:   d[1],
		}
	}
	conf = Config{
		sid:        jsonObj.Sid,
		pwd:        jsonObj.Password,
		allBooking: allBooking,
	}
	conf.vip = jsonObj.Vip
	if !checkConf(conf) {
		log.Fatalln("配置文件存在问题。")
	}
	return
}
