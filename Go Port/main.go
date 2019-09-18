package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	urlLogin   = "http://202.120.82.2:8081/ClientWeb/pro/ajax/login.aspx"
	urlBooking = "http://202.120.82.2:8081/ClientWeb/pro/ajax/reserve.aspx"
)

var (
	client = &http.Client{}
)

func login(sid string, pwd string) (cookie string, ans bool) {
	resp, err := http.PostForm(urlLogin,
		url.Values{
			"id":  {sid},
			"pwd": {pwd},
			"act": {"login"}})
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	cookie = ""
	if strings.Contains(string(body), "\"msg\":\"ok\"") {
		ans = true
		for _, v := range resp.Cookies() {
			if v.Name == "ASP.NET_SessionId" {
				cookie = v.Value
				return
			}
		}
	} else {
		ans = false
		cookie = string(body)
	}
	return
}

func booking(req *http.Request, c chan string) {
	resp, err := client.Do(req)
	if err != nil {
		// 这里必须处理错误，因为会超时。
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	strBody := string(body)
	log.Println(strBody)
	if !strings.Contains(strBody, "要到[21:00]方可预约") {
		c <- strBody
	}
}

func getBookingReq(room Room, startTime string, endTime string, delayDay int) *http.Request {
	req, err := http.NewRequest("GET", urlBooking, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 后天日期, 2019-09-19
	theDayAfterT := time.Now().AddDate(0, 0, delayDay).Format("2006-01-02")

	q1 := req.URL.Query()
	q1.Add("dialogid", "")
	q1.Add("dev_id", room.devId)
	q1.Add("lab_id", room.labId)
	q1.Add("kind_id", room.kindId)
	q1.Add("room_id", "")
	q1.Add("type", "dev")
	q1.Add("prop", "")
	q1.Add("test_id", "")
	q1.Add("term", "")
	q1.Add("test_name", "")
	q1.Add("up_file", "")
	q1.Add("memo", "")
	q1.Add("act", "set_resv")
	//q1.Add("_", "")
	q1.Add("start", theDayAfterT+" "+startTime)
	q1.Add("end", theDayAfterT+" "+endTime)
	q1.Add("start_time", "")
	q1.Add("end_time", "")
	req.URL.RawQuery = q1.Encode()
	return req
}

func main() {
	log.SetFlags(log.Lmicroseconds | log.LstdFlags)
	log.Println("================= Start ===================")

	conf := GetConf()
	// 1. 通过登录获取带认证的cookie
	sessionId, succeed := login(conf.sid, conf.pwd)
	if !succeed {
		log.Println("登录失败！")
		log.Println(sessionId)
		return
	}
	log.Println("登录成功")

	// 2. 通过带认证的 cookie 构造 带参数的request请求
	theB := conf.allBooking[0]
	req1 := getBookingReq(theB.room, theB.startTime, theB.endTime, theB.delayDay)
	theB = conf.allBooking[1]
	req2 := getBookingReq(theB.room, theB.startTime, theB.endTime, theB.delayDay)
	theB = conf.allBooking[2]
	req3 := getBookingReq(theB.room, theB.startTime, theB.endTime, theB.delayDay)
	// set cookie
	cookie := http.Cookie{Name: "ASP.NET_SessionId", Value: sessionId}
	req1.AddCookie(&cookie)
	req2.AddCookie(&cookie)
	req3.AddCookie(&cookie)

	// 3. 疯狂发送构造好的request
	c1 := make(chan string)
	c2 := make(chan string)
	c3 := make(chan string)
	finish1 := false
	finish2 := false
	finish3 := false
	succeed1 := false
	succeed2 := false
	succeed3 := false

	for (!finish1) || (!finish2) || (!finish3) {
		if !finish1 {
			go booking(req1, c1)
		}
		if !finish2 {
			go booking(req2, c2)
		}
		if !finish3 {
			go booking(req3, c3)
		}
		select {
		case tmp := <-c1:
			if strings.Contains(tmp, "操作成功") {
				succeed1 = true
			}
			finish1 = true
		case tmp := <-c2:
			if strings.Contains(tmp, "操作成功") {
				succeed2 = true
			}
			finish2 = true
		case tmp := <-c3:
			if strings.Contains(tmp, "操作成功") {
				succeed3 = true
			}
			finish3 = true
		default:
			time.Sleep(30 * time.Millisecond)
		}
	}
	theB = conf.allBooking[0]
	if succeed1 {
		log.Println(theB.startTime + " -- " + theB.endTime + "  " + theB.room.devName + "----预定成功----")
	}
	theB = conf.allBooking[1]
	if succeed2 {
		log.Println(theB.startTime + " -- " + theB.endTime + "  " + theB.room.devName + "----预定成功----")
	}
	theB = conf.allBooking[2]
	if succeed3 {
		log.Println(theB.startTime + " -- " + theB.endTime + "  " + theB.room.devName + "----预定成功----")
	}

	log.Println("================= Over ===================")
}
