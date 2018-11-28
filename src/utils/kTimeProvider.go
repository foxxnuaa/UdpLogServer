package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var G_NowTime time.Time
var G_NowTimeLock sync.RWMutex

func StartNowTime() {
	go func() {
		for {
			G_NowTimeLock.Lock()
			G_NowTime = time.Now()
			G_NowTimeLock.Unlock()
			time.Sleep(time.Duration(1) * time.Second)
		}

	}()
}

func GetNowTime() time.Time {
	G_NowTimeLock.RLock()
	defer G_NowTimeLock.RUnlock()
	NowTime := G_NowTime
	if NowTime.IsZero() {
		return time.Now()
	}
	return NowTime

}

func GetDate() int {
	NowTime := GetNowTime()
	strDate := NowTime.Format("20060102")
	iDate, _ := strconv.Atoi(strDate)
	return iDate
}

func GetDateBefore(BeforeDay int) int {
	NowTime := GetNowTime()
	BeforeTime := NowTime.AddDate(0, 0, -1*BeforeDay)
	strDate := BeforeTime.Format("20060102")
	iDate, _ := strconv.Atoi(strDate)
	return iDate
}

func GetSecondFromDate(iData int) uint32 {
	//func Date(year int, month Month, day, hour, min, sec, nsec int, loc *Location) Time
	iYear := iData / 10000
	var iMonth time.Month
	iMonth = time.Month((iData - iYear*10000) / 100)
	iDay := (iData - iYear*10000 - int(iMonth)*100)
	stTime := time.Date(iYear, iMonth, iDay, 0, 0, 0, 0, time.Local)
	return uint32(stTime.Unix())
}

func GetDateByUnixTime(unixTime int64) int {

	thatTime := time.Unix(unixTime, 0)
	strDate := thatTime.Format("20060102")
	iDate, _ := strconv.Atoi(strDate)
	return iDate
}

func IsToday(unixTime int64) bool {
	y, m, d := GetNowTime().Date()
	ty, tm, td := time.Unix(unixTime, 0).Date()
	if y == ty && tm == m && td == d {
		return true
	}
	return false

}

func GetParseInLocationTime(year int, month time.Month, day, hour, min, second int) time.Time {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	var strShortString string
	var strAmPm string = "am"
	nowTime := time.Now()
	strZoneName, _ := nowTime.Zone()
	if hour >= 12 {
		hour = hour - 12
		strAmPm = "pm"
	}
	strShortString = string(month.String()[0:3])
	var strMin string = strconv.Itoa(min)
	if min < 10 {
		strMin = "0" + strMin
	}
	targetTimeString := fmt.Sprintf("%s %d, %d at %d:%s%s (%s)", strShortString, day, year, hour, strMin, strAmPm, strZoneName)
	targetTime, err := time.ParseInLocation(longForm, targetTimeString, time.Local)
	if err != nil {
		os.Exit(0)
	}
	return targetTime
}

func GetNext5RefreshTime() int64 {
	now := GetNowTime()
	var refreshTime int64
	if now.Hour() < 5 {
		refreshTime = GetParseInLocationTime(now.Year(), now.Month(), now.Day(), 5, 0, 0).Unix()
	} else {
		tomorrow := now.Add(24 * time.Hour)
		refreshTime = GetParseInLocationTime(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 5, 0, 0).Unix()
	}
	return refreshTime
}

func Get5RefreshTime() int64 {
	now := GetNowTime()
	var refreshTime int64
	refreshTime = GetParseInLocationTime(now.Year(), now.Month(), now.Day(), 5, 0, 0).Unix()

	return refreshTime
}

// 将普通日期转化是unix时间戳，2015-10-15 14:20:33
func NormalDateToUnixTime(strDate string) (openTime int64) {
	const longForm = "2006-01-02 at 3:04pm (MST)"

	vecYm := strings.Split(strDate, " ")
	if len(vecYm) < 2 {
		fmt.Println("NormalDateToUnixTime split strDate fail", strDate)
		return 0
	}
	strYear := vecYm[0]
	strTime := vecYm[1]

	vecTime := strings.Split(strTime, ":")
	if len(vecYm) < 2 {
		fmt.Println("NormalDateToUnixTime split strTime fail", strTime)
		return 0
	}
	strHour := vecTime[0]
	strMinute := vecTime[1]

	hour, _ := strconv.Atoi(strHour)
	var strAmPm string = "am"
	nowTime := time.Now()
	strZoneName, _ := nowTime.Zone()
	if hour >= 12 {
		hour = hour - 12
		strAmPm = "pm"
	}
	strHour = strconv.Itoa(hour)

	targetTimeString := fmt.Sprintf("%s at %s:%s%s (%s)", strYear, strHour, strMinute, strAmPm, strZoneName)

	targetTime, err := time.ParseInLocation(longForm, targetTimeString, time.Local)
	if err != nil {
		fmt.Println("ParseInLocation err", err)
		return 0
	}

	openTime = targetTime.Unix()
	return openTime
}
