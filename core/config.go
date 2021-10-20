package core

// HTTP(S) --> http(s)://example.com:8080
// TCP Bind --> tcp://127.0.0.1:5555, tcp://0.0.0.0:5555
// TCP Reverse --> tcp://192.168.1.100:8888

var (
	Sleep   = 60000
	Jitter  = 10
	C2_URL  = "http://127.0.0.1:80"
	GetURL  = C2_URL + "/load"
	PostURL = C2_URL + "/submit.php?id=%d"
	Headers = map[string]string{
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:69.0) Gecko/20100101 Firefox/69.0",
	}
	Public_Key = "-----BEGIN PUBLIC KEY-----\nXXXYYYZZZ\n-----END PUBLIC KEY-----"
)
