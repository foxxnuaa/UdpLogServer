package main

import (
	"os"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"utils"
	"strings"
	"net"
	"log"
	"io"
	"fmt"
	"time"
	"reflect"
	"errors"
	"sort"
	"html/template"
	"bufio"
	"sync"
	"strconv"
	"math"
)

type ServerConf struct {
	MysqlConf string
	HttpPort  string
	UdpPort   string
}

type FileItem struct {
	FileName       string
	FileDir        string
	FileModifyTime int64
	FileCreateTime string
}

type FileContent struct {
	Level   string
	Content string
}

var G_StConf ServerConf
var FileCreateTimeMap sync.Map

func loadConf() {
	fi, err := os.Open("conf.json")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	//utils.Debugln("loadConf:", string(fd))
	err = json.Unmarshal(fd, &G_StConf)
	if err != nil {
		utils.Debugln("error:", err)
		os.Exit(0)
	}
	return
}

func FunctionMapCall(m map[string]interface{}, name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(m[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

func ajaxHandler(w http.ResponseWriter, r *http.Request) {
	utils.Debugln("ajaxHandler begin---------------")
	funcs := make(map[string]interface{})
	funcs["getTemplates"] = getTemplates
	funcs["getdirs"] = getdirs
	funcs["getFiles"] = getFiles
	funcs["getFileContentWithLevel"] = getFileContentWithLevel
	funcs["getFileContent"] = getFileContent
	funcs["getRoomFiles"] = getRoomFiles

	strAction := r.FormValue("action")
	strData := r.FormValue("data")
	utils.Debugln("ajaxHandler---------------", strAction, strData)
	FunctionMapCall(funcs, strAction, w, r, strData)
}

func GetRetMap() map[string]string {
	var retDataMap map[string]string = make(map[string]string)
	retDataMap["ret"] = "0"
	retDataMap["data"] = ""
	return retDataMap
}

func RetMap2String(retMap map[string]string) string {
	v, _ := json.Marshal(retMap)
	return string(v)
}

func getFileContent(w http.ResponseWriter, r *http.Request, strData string) {
	fmt.Println("getFileContent begin")
	retMap := GetRetMap()
	pstFile, err := os.OpenFile("log/"+strData, os.O_RDONLY, 0666)
	if err != nil {
		pstFile, err = os.OpenFile("room/"+strData, os.O_RDONLY, 0666)
	}
	var FileContentList []FileContent
	if (err == nil) {
		buf := bufio.NewReader(pstFile)
		for {
			line, err := buf.ReadString('\n')
			vStrings := strings.Split(line, "|")
			if (len(vStrings) < 3) {
				break
			}
			var stFileContent FileContent
			stFileContent.Level = vStrings[1]
			stFileContent.Content = line
			FileContentList = append(FileContentList, stFileContent)
			if err != nil {
				if err == io.EOF {
					break
				}
			}
		}
		pstFile.Close()
	}
	Buffer, err := json.Marshal(FileContentList)
	if (err == nil) {
		retMap["data"] = string(Buffer)
	}
	fmt.Fprint(w, RetMap2String(retMap))
}

func getFileContentWithLevel(w http.ResponseWriter, r *http.Request, strData string) {
	fmt.Println("getFileContent begin")
	var ParamMap  map[string]string
	json.Unmarshal([]byte(strData),&ParamMap)

    PangeNo,_ := strconv.Atoi(ParamMap["PageNum"])
	PageSize,_ := strconv.Atoi(ParamMap["PageSize"])
    nFromIndex := (PangeNo-1) * PageSize
    nEndIndex := PangeNo * PageSize
    fmt.Println("nFromIndex==",nFromIndex,"nEndIndex==",nEndIndex)
	vParam := strings.Split(ParamMap["Levels"], "|")
	var levelMap  = make(map[string]bool)
	if len(vParam) == 0 {
		levelMap["0"] = true
		levelMap["1"] = true
		levelMap["2"] = true
		levelMap["3"] = true
		levelMap["4"] = true
		levelMap["5"] = true
	} else {
		for i := 0; i < len(vParam); i++ {
			levelMap[vParam[i]] = true
		}
	}
	fmt.Println("levels=", levelMap)
	retMap := GetRetMap()
	PthSep := string(os.PathSeparator)
	pstFile, err := os.OpenFile("log"+PthSep+ParamMap["FileName"], os.O_RDONLY, 0666)
	if err != nil {
		pstFile, err = os.OpenFile("room"+PthSep+ParamMap["FileName"], os.O_RDONLY, 0666)
	}
	TotalPageNum := 0
	var FileContentList []FileContent
	if err == nil {
		buf := bufio.NewReader(pstFile)
		nLineIndex:=0
		for {

			line, err := buf.ReadString('\n')
			//fmt.Println("line=", line)
			vStrings := strings.Split(line, "|")
			if len(vStrings) < 3{
				break
			}

			var stFileContent FileContent
			stFileContent.Level = vStrings[1]
			stFileContent.Content = line
			if levelMap[stFileContent.Level] {
				if nLineIndex >=nFromIndex && nLineIndex <= nEndIndex{
					FileContentList = append(FileContentList, stFileContent)
				}
				nLineIndex++

			}

			if err != nil {
				if err == io.EOF {
					break
				}
			}
		}
		fmt.Println("nLineIndex==",nLineIndex)
		TotalPageNum = int(math.Ceil(float64(nLineIndex) / float64(PageSize)))
		pstFile.Close()
	}
	Buffer, err := json.Marshal(FileContentList)
	PageInfoAndFileContent := make(map[string]string)
	PageInfoAndFileContent["FileContent"] = string(Buffer)
	PageInfoAndFileContent["TotalPageNum"] = strconv.Itoa(TotalPageNum)
	if err == nil {
		strData,_ := json.Marshal(PageInfoAndFileContent)
		retMap["data"] = string(strData)
	}
	fmt.Fprint(w, RetMap2String(retMap))
}

func getdirs(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	Files, err := ListDir("log", "")
	if (err == nil) {
		buff, _ := json.Marshal(Files)
		retMap["data"] = string(buff)
	}
	fmt.Fprint(w, RetMap2String(retMap))
}

func getFiles(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	Files, err := ListFile("log/"+strData, "")
	if (err == nil) {
		buff, _ := json.Marshal(Files)
		retMap["data"] = string(buff)
	}
	fmt.Fprint(w, RetMap2String(retMap))
}

func getRoomFiles(w http.ResponseWriter, r *http.Request, strData string) {
	retMap := GetRetMap()
	Files, err := ListFile("room", "")
	if (err == nil) {
		buff, _ := json.Marshal(Files)
		retMap["data"] = string(buff)
	}
	fmt.Fprint(w, RetMap2String(retMap))
}

func getTemplates(w http.ResponseWriter, r *http.Request, strData string) {
	var FileList []string
	json.Unmarshal([]byte(strData), &FileList)
	var templatesMap map[string]string = make(map[string]string)
	retMap := GetRetMap()
	for i := 0; i < len(FileList); i++ {
		strFullName := "template/" + FileList[i]
		fi, err := os.Open(strFullName)
		if err != nil {
			panic(err)
		}
		defer fi.Close()
		fd, err := ioutil.ReadAll(fi)
		templatesMap[FileList[i]] = string(fd)
		utils.Debugln("getTemplates:", string(fd))

	}
	templatesFiles2String := RetMap2String(templatesMap)
	retMap["data"] = templatesFiles2String
	fmt.Fprint(w, RetMap2String(retMap))
}

func StartWeb() error {
	//http.Handle("/", http.FileServer(http.Dir("./log")))
	http.Handle("/1/", http.StripPrefix("/1/", http.FileServer(http.Dir("./"))))
	http.Handle("/", http.FileServer(http.Dir("resources/")))
	http.HandleFunc("/index", leftbarHandler)
	http.HandleFunc("/ajax", ajaxHandler)
	err := http.ListenAndServe(G_StConf.HttpPort, nil)
	return err

}

func GetFileModTime(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		log.Println("open file error")
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Println("stat fileinfo error")
		return time.Now().Unix()
	}
	return fi.ModTime().Unix()
}

func GetFileCreateTime(path string) string {
	strRet, OK := FileCreateTimeMap.Load(path)
	if (OK) {
		return strRet.(string)
	}

	CreateTime:=time.Unix(GetFileModTime(path),0)
	return CreateTime.Format("2006/01/02 15:04:05")
}

func leftbarHandler(w http.ResponseWriter, r *http.Request) {

	//ServerGroups := GetServerGroups()
	t, err := template.ParseFiles("template/leftbar.html")
	if err != nil {
		utils.Debugln("ParseFiles error", err.Error())
	}

	err = t.Execute(w, nil)
	if err != nil {
		utils.Debugln("ParseFiles error", err.Error())
	}

}

type FileItems []FileItem

func (a FileItems) Len() int           { return len(a) }
func (a FileItems) Less(i, j int) bool { return (a[i].FileModifyTime < a[j].FileModifyTime) }
func (a FileItems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func ListDir(dirPth string, suffix string) (files FileItems, err error) {
	//files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if !fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			var stFileItem FileItem
			stFileItem.FileModifyTime = GetFileModTime(dirPth + PthSep + fi.Name())
			stFileItem.FileName = fi.Name()
			files = append(files, stFileItem)
		}
	}
	sort.Sort(files)
	return files, nil
}

func ListFile(dirPth string, suffix string) (files FileItems, err error) {
	//files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			var stFileItem FileItem
			stFileItem.FileModifyTime = GetFileModTime(dirPth + PthSep + fi.Name())
			stFileItem.FileName = fi.Name()
			stFileItem.FileDir = dirPth[4:]
			stFileItem.FileCreateTime = GetFileCreateTime(dirPth + PthSep + fi.Name())
			files = append(files, stFileItem)
		}
	}
	sort.Sort(files)
	return files, nil
}

func DealWithUdpPkg() {
	//f, err1 := os.OpenFile("alian", os.O_APPEND, 0777)
	var FileMap map[string]*os.File = make(map[string]*os.File);
	udpAddr, err := net.ResolveUDPAddr("udp", G_StConf.UdpPort)
	if err != nil {
		log.Fatalln("Error: ", err)
		os.Exit(0)
	}

	// Build listining connections
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalln("Error: ", err)
		os.Exit(0)
	}
	defer conn.Close()
	PthSep := string(os.PathSeparator)
	// Interacting with one client at a time
	recvBuff := make([]byte, 1500)
	for {
		//log.Println("Ready to receive packets!")
		// Receiving a message
		rn, clientaddr, err := conn.ReadFromUDP(recvBuff)
		if rn == 1 {
			conn.WriteToUDP([]byte(recvBuff[0:1]), clientaddr)
		}
		if err != nil {
			utils.Errorln("Error:", err.Error())
			time.Sleep(time.Millisecond * 30)
			continue
		}
		//utils.Debugln("DealWithUdpPkg",string(recvBuff))
		recvBuff2 := recvBuff[:rn]
		var strReceived string = string(recvBuff2)
		strReceived = strings.Replace(strReceived, "\n", "**&", -1)
		strReceived = strings.Replace(strReceived, "\r", "**&", -1)
		//fmt.Println(strReceived)
		vLogItem := strings.Split(strReceived, "|")
		vRoomLogItem := strings.Split(strReceived, "]")
		if (len(vLogItem) < 3) {
			continue
		}
		strFileDate := time.Now().Format("2006-01-02")
		nHour, nMinute, nSecond := time.Now().Clock()
		strNowTime := fmt.Sprintf("%s:%2.2d:%2.2d:%2.2d|", strFileDate, nHour, nMinute, nSecond)
		_, err = os.Stat("log"+PthSep + strFileDate)
		if err != nil {
			os.Mkdir("log"+PthSep+strFileDate, 0777)
		}
		strFileName := "log"+PthSep + strFileDate +PthSep + fmt.Sprintf(strFileDate+"_"+vLogItem[1]) + ".log"
		if pfile, OK := FileMap[strFileName]; OK {
			io.WriteString(pfile, strNowTime)
			io.WriteString(pfile, strReceived)
			io.WriteString(pfile, "\r\n")
		} else {
			pfile, _ = os.OpenFile(strFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			_, OK := FileCreateTimeMap.Load(strFileName)
			if !OK {
				FileCreateTimeMap.Store(strFileName, time.Now().Format("2006/01/02 15:04:05"))
			}
			FileMap[strFileName] = pfile
			io.WriteString(pfile, strNowTime)
			io.WriteString(pfile, strReceived)
			io.WriteString(pfile, "\r\n")
		}

		//写room日志
		if len(vRoomLogItem) < 3 {
			continue
		}
		strRoom := vRoomLogItem[1][2:]
		if (strRoom == "0") {
			continue
		}
		_, err = os.Stat("room")
		if err != nil {
			os.Mkdir("room", 0777)
		}
		strRoomFileName := "room/" + strRoom + ".log"
		if pfile2, OK := FileMap[strRoomFileName]; OK {
			io.WriteString(pfile2, strNowTime)
			io.WriteString(pfile2, strReceived)
			io.WriteString(pfile2, "\r\n")
		} else {
			pfile2, _ = os.OpenFile(strRoomFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			FileMap[strRoomFileName] = pfile2
			io.WriteString(pfile2, strNowTime)
			io.WriteString(pfile2, strReceived)
			io.WriteString(pfile2, "\r\n")
		}

	}
}
func main() {
	//StartMysqlWriteTimer()
	utils.SetLevel(4)
	loadConf()
	go DealWithUdpPkg()
	time.Sleep(time.Second * 3)
	conn, err := net.Dial("udp", "127.0.0.1:10433")
	defer conn.Close()
	if err != nil {
		os.Exit(1)
	}
	conn.Write([]byte("1|阿涟|[00:31][10991] Hello\n \r world!"))
	conn.Write([]byte("0"))
	var rdbuff []byte
	conn.Read(rdbuff)
	fmt.Println(string(rdbuff))
	StartWeb()
}
