package utils

import (
	"net/http"
	"net/url"
	"fmt"
	"io/ioutil"
	"encoding/json"
)
type AgentConfig struct {
	StrOuterIp       string
	StrInnerIp     string
	StrWorkDir string
	StrHttpPort string
	StrCenterUrl string
}

var G_StAgentConf AgentConfig
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
func HttpPostForm(strUrl string, data *map[string]string) (bRetCode bool, strOut string) {

	v := url.Values{}
	for key, value := range *data {
		v.Set(key, value)
	}
	fmt.Println("v.Encode()=", v.Encode())
	resp, err := http.PostForm(strUrl, v)
	if err != nil {
		fmt.Println("HttpPostForm failed!", strUrl, v)
		return false, ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("HttpPostForm failed!", strUrl, data)
		return false, ""
	}
	strOut = string(body)

	return true, strOut
}
