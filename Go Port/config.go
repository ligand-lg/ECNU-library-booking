package main

type Room struct {
	devId   string
	kindId  string
	labId   string
	devName string
	roomNo  string
}

type Booking struct {
	room      Room
	delayDay  int
	startTime string
	endTime   string
}

var (
	C424 = Room{
		devId:   "3676522",
		devName: "中北校区单人间C424",
		kindId:  "3675133",
		roomNo:  "C424",
		labId:   "3674920",
	}
	C426 = Room{
		devId:   "3676547",
		devName: "中北校区单人间C426",
		kindId:  "3675133",
		roomNo:  "C426",
		labId:   "3674920",
	}
)

type Config struct {
	sid        string
	pwd        string
	allBooking [3]Booking
}

func GetConf() (conf Config) {
	//content, err := ioutil.ReadFile("conf.json")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//var tmp interface{}
	//err = json.Unmarshal(content, tmp)
	//if err != nil {
	//	log.Fatal("配置文件不是合法的JSON")
	allBooking := [3]Booking{
		Booking{
			room:      C424,
			delayDay:  2,
			startTime: "09:30",
			endTime:   "13:30",
		}, Booking{
			room:      C424,
			delayDay:  2,
			startTime: "13:50",
			endTime:   "17:50",
		}, Booking{
			room:      C424,
			delayDay:  2,
			startTime: "18:00",
			endTime:   "22:00",
		}}
	conf = Config{
		sid:        "",
		pwd:        "",
		allBooking: allBooking }

	return
}
