package core

import (
	"bytes"
	"encoding/binary"
	"io"
	"math/rand"
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

func JoinMap(Map ...map[string]string) map[string]string {
	New := map[string]string{}
	for _, x := range Map {
		for Key, Value := range x {
			New[Key] = Value
		}
	}
	return New
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
