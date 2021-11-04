package core

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	TRUE = true
	MLock  sync.Mutex
	RWLock sync.RWMutex
	Buffer bytes.Buffer
)

func IntToByte(Len, Num int) []byte {
	Buf := make([]byte, Len)
	switch Len {
	case 1:
		Buf[0] = byte(Num)
	case 2:
		binary.BigEndian.PutUint16(Buf, uint16(Num))
	case 4:
		binary.BigEndian.PutUint32(Buf, uint32(Num))
	}
	return Buf
}

func _IntToByte(Len, Num int) []byte {
	Buf := make([]byte, Len)
	switch Len {
	case 1:
		Buf[0] = byte(Num)
	case 2:
		binary.LittleEndian.PutUint16(Buf, uint16(Num))
	case 4:
		binary.LittleEndian.PutUint32(Buf, uint32(Num))
	}
	return Buf
}

func ByteToInt(Bytes []byte) int {
	switch len(Bytes) {
	case 1:
		return int(Bytes[0])
	case 2:
		return int(binary.BigEndian.Uint16(Bytes))
	case 4:
		return int(binary.BigEndian.Uint32(Bytes))
	}
	return 0
}

func _ByteToInt(Bytes []byte) int {
	switch len(Bytes) {
	case 1:
		return int(Bytes[0])
	case 2:
		return int(binary.LittleEndian.Uint16(Bytes))
	case 4:
		return int(binary.LittleEndian.Uint32(Bytes))
	}
	return 0
}

func BigToLittle(Bytes []byte) []byte {
	return _IntToByte(len(Bytes), ByteToInt(Bytes))
}

func LittleToBig(Bytes []byte) []byte {
	return IntToByte(len(Bytes), _ByteToInt(Bytes))
}

func ToLength(s string) int {
	v, err := strconv.Atoi(s)
	if err == nil { return v }
	return len(s)
}

func JoinBytes(Bytes ...[]byte) []byte {
	return bytes.Join(Bytes, []byte(""))
}

func RandomInt(Min, Max int) int {
	rand.Seed(time.Now().UnixNano())
	if Max > Min {
		return Min + rand.Intn(Max-Min)
	}
	if Max < Min {
		return Max + rand.Intn(Min-Max)
	}
	return Max
}

func ReadFrom(r io.Reader, Len int) []byte {
	Buf := make([]byte, Len)
	_, err := io.ReadFull(r, Buf)
	if err != nil { return nil }
	return Buf
}

func MutexLock(Worker func()) {
	MLock.Lock()
	Worker()
	MLock.Unlock()
}

func WriteLock(Worker func())  {
	RWLock.Lock()
	Worker()
	RWLock.Unlock()
}

func Base64ToMap(Base64 string) map[string]string {
	Map := map[string]string{}
	Bytes, _ := base64.StdEncoding.DecodeString(Base64)
	for _, v := range strings.Split(string(Bytes), "\n") {
		if !strings.Contains(v, ":") { continue }
		Line := strings.SplitN(v, ":", 2)
		Key  := strings.TrimSpace(Line[0])
		Val  := strings.TrimSpace(Line[1])
		if Key != "" || Val != "" {
			Map[Key] = Val
		}
	}
	return Map
}