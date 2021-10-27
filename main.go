package main

import "Geacon/core"

func main() {
	for core.TRUE {
		BUF := core.Pull()
		for BUF != nil && BUF.Len() >= 8 {
			ID := core.ByteToInt(BUF.Next(4))
			Data := BUF.Next(core.ByteToInt(BUF.Next(4)))
			if FUNC, OK := core.TK[ID]; OK {
				go FUNC(Data)
			}
		}
		core.Push()
	}
}