// 游戏主程序 - 通信
package Game

import (
	"TestMMO2/src/Game/PB"
	"TestMMO2/src/MServer"
	"bytes"
	"fmt"
)

type CallBackRouter func(session *MServer.MSession, bitmsg []byte)

// 初始化路由
func (g *Game) InitRouter() {
	g.RouterDict = make(map[uint32]CallBackRouter)

	g.AddRouter((uint32)(PB.MSG_ID_LOGIN), g.C2S_LOGIN)
	g.AddRouter((uint32)(PB.MSG_ID_MOVE_TO), g.C2S_MOVETO)
	g.AddRouter((uint32)(PB.MSG_ID_MOVE_STOP), g.C2S_MOVESTOP)

	// g.RouterDict[ENUM_C2S_LOGIN] = g.C2S_LOGIN
	// g.RouterDict[ENUM_MOVETO] = g.C2S_MOVETO
	// g.RouterDict[ENUM_STOPMOVE] = g.C2S_STOPMOVE
}

// 添加路由
func (g *Game) AddRouter(msgid uint32, router CallBackRouter) {
	g.RouterDict[msgid] = router
}

// 删除路由
func (g *Game) RemoveRouter() {

}

// 新链接
func (g *Game) SessionAccpet(s *MServer.MSession) {
	fmt.Println("连接", s.ADDR)
}

// 链接错误
func (g *Game) SessionError(s *MServer.MSession, err error) {

	g.PlayerOffline(s.ID)
	//fmt.Println("断开", s.ID, s.ADDR, "  ", err)
}

// 接收消息
//  4个字节的ID + 4个字节的长度 + Message Byte
func (g *Game) SessionRecv(s *MServer.MSession, msgid uint32, bit_msg []byte) {

	//fmt.Println("Recv ", s.ID, msgid)
	callfun, ok := g.RouterDict[msgid]
	if ok {
		callfun(s, bit_msg)
	}

}

// 准备发送
func (g *Game) ReadySend(msgid uint32, s *MServer.MSession, bitmsg []byte) {
	s.Send(msgid, bitmsg)
}

// 广播
func (g *Game) BroadCast(msgid uint32, bitmsg []byte) {
	g.server.BroadCast(msgid, bitmsg)
}

//BytesCombine
func BytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}
