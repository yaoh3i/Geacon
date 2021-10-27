package core

import (
	"net/url"
	"strings"
)

func GetURL(Map map[string]interface{}) string {
	return Map["URL"].(string)
}

func GetHeader(Map map[string]interface{}) map[string]string {
	return Map["Header"].(map[string]string)
}

func GetOutput(Bytes []byte) []byte {
	Info := PullInfo["Output"].(map[string]string)
	Data := TrimPSfix(string(Bytes), Info["Prepend"], Info["Append"])
	return Decoding(Data, Info["Coding"])
}

func SetOutput(Bytes []byte) *strings.Reader {
	Info := PushInfo["Output"].(map[string]string)
	Data := Info["Prepend"] + Encoding(Bytes, Info["Coding"]) + Info["Append"]
	return strings.NewReader(Data)
}

func SetProperty(Map map[string]interface{}, Key string, Value []byte) {
	URL, _ := url.Parse(GetURL(Map))
	Header := Base64ToMap(Map["Header"].(string))
	KeyMap := Map[Key].(map[string]string)
	KeyInfo := strings.SplitN(KeyMap["Store"], ":", 2)
	KeyData := KeyMap["Prepend"] + Encoding(Value, KeyMap["Coding"]) + KeyMap["Append"]
	if strings.ToUpper(KeyInfo[0]) == "URL" {
		if len(KeyInfo) == 2 && len(KeyInfo[1]) > 0 {
			Query := URL.Query()
			Query.Set(KeyInfo[1], KeyData)
			URL.RawQuery = Query.Encode()
		} else {
			URL.Path = URL.Path + KeyData
		}
	}
	if strings.ToUpper(KeyInfo[0]) == "HEADER" {
		Header[KeyInfo[1]] = KeyData
	}
	Map["URL"] = URL.String()
	Map["Header"] = Header
}