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
	DEBUG         = false
	urlLogin      = "http://202.120.82.2:8081/ClientWeb/pro/ajax/login.aspx"
	urlBooking    = "http://202.120.82.2:8081/ClientWeb/pro/ajax/reserve.aspx"
	sessionIdName = "ASP.NET_SessionId"
)

var (
	client = &http.Client{}
)

const (
	statusNotFinish = 0
	statusSucceed   = 1
	statusFailed    = 2
)

func encodeResult(id, status int) int {
	return id*10 + status
}
func decodeResult(res int) (id int, status int) {
	id = res / 10
	status = res % 10
	return
}

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
	if strings.Contains(string(body), "\"msg\":\"ok\"") {
		ans = true
		for _, v := range resp.Cookies() {
			if v.Name == sessionIdName {
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

func booking(req *http.Request, id int, c chan int) {
	resp, err := client.Do(req)
	if err != nil {
		// 这里必须处理错误，因为之后会使用到resp，如果其为nil，程序会崩掉。
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	strBody := string(body)
	// 这里必须打印日志，否则后台一直疯狂post数据，你却不知道。
	log.Println(strBody)

	// 除时间未到的情况外，其他情况都直接终止当前任务
	if !strings.Contains(strBody, "要到[21:00]方可预约") {
		if strings.Contains(strBody, "操作成功") {
			c <- encodeResult(id, statusSucceed)
		} else {
			c <- encodeResult(id, statusFailed)
		}
	}
}

// 设置预定字段
func getBookingReq(booking Booking) *http.Request {
	req, err := http.NewRequest("GET", urlBooking, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 后天日期, 2019-09-19
	theDayAfterT := time.Now().AddDate(0, 0, booking.delayDay).Format("2006-01-02")

	q1 := req.URL.Query()
	q1.Add("dialogid", "")
	q1.Add("dev_id", booking.room.DevId)
	q1.Add("lab_id", booking.room.LabId)
	q1.Add("kind_id", booking.room.KindId)
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
	q1.Add("start", theDayAfterT+" "+booking.startTime)
	q1.Add("end", theDayAfterT+" "+booking.endTime)
	q1.Add("start_time", "")
	q1.Add("end_time", "")
	req.URL.RawQuery = q1.Encode()
	return req
}

/* 对比当前时间与9点的间隔
vip: 20:59:50
非vip: 20:59:55
*/

func checkTime(isVip bool) {
	desStr := "20:59:55"
	if isVip {
		desStr = "20:59:50"
	}
	layout := "15:04:05"
	nowStr := time.Now().Format(layout)
	des, _ := time.Parse(layout, desStr)
	now, _ := time.Parse(layout, nowStr)
	diff := des.Sub(now)
	if diff > 0*time.Second {
		log.Println("等待：", diff)
		time.Sleep(diff)
	}
}

func main() {
	// 设置时间格式为微秒级别，用于分析时延
	log.SetFlags(log.Lmicroseconds | log.LstdFlags)
	log.Println("================= Start ===================")

	conf := GetConf()
	if !DEBUG {
		checkTime(conf.vip)
	}

	// 1. 通过登录获取带认证的cookie
	sessionId, loginSucceed := login(conf.sid, conf.pwd)
	if !loginSucceed {
		log.Println("登录失败！")
		log.Println(sessionId)
		return
	}
	log.Println("登录成功")

	// 2. 通过带认证的 cookie 构造带参数的request请求
	cookie := http.Cookie{Name: sessionIdName, Value: sessionId}

	bookingCnt := len(conf.allBooking)
	reqs := make([]*http.Request, bookingCnt)
	for i, booking := range conf.allBooking {
		reqs[i] = getBookingReq(booking)
		reqs[i].AddCookie(&cookie)
	}
	// buffered channels
	c := make(chan int, bookingCnt)
	allStatus := make([]int, bookingCnt)
	for i, _ := range allStatus {
		allStatus[i] = statusNotFinish
	}

	for allFinish := false; !allFinish; {
		suicide()
		allFinish = true
		//
		for i, status := range allStatus {
			if status == statusNotFinish {
				allFinish = false
				go booking(reqs[i], i, c)
			}
		}
		for flag := true; flag; {
			select {
			case res := <-c:
				id, status := decodeResult(res)
				if allStatus[id] != statusSucceed {
					allStatus[id] = status
				}
			default:
				flag = false
			}
		}
		time.Sleep(30 * time.Millisecond)
	}
	// 3. 输出结果
	for i, b := range conf.allBooking {
		switch allStatus[i] {
		case statusSucceed:
			log.Println(b.startTime + " -- " + b.endTime + "  " + b.room.DevName + "   ----预定成功----")
		case statusFailed:
			log.Println(b.startTime + " -- " + b.endTime + "  " + b.room.DevName + "   xxxx预定失败xxxx")
		default:
			log.Println("未知状态：", allStatus[i])
		}
	}
	log.Println("================= Over ===================")
}
