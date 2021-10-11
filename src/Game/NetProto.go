//
package Game

import (
	"TestMMO2/src/Game/PB"
	"TestMMO2/src/MServer"

	"fmt"

	proto "github.com/gogo/protobuf/proto"
)

// 登录
func (g *Game) C2S_LOGIN(s *MServer.MSession, bitmsg []byte) {
	mLogin := PB.Login{}
	proto.Unmarshal(bitmsg, &mLogin)

	fmt.Println("登陆 : ", s.ADDR, mLogin.Id, mLogin.Pwd)

	bitsysinfo, _ := proto.Marshal(g.GetSysInfo())
	s.Send((uint32)(PB.MSG_ID_SYS_INFO), bitsysinfo) // 发送系统信息

	mUser := g.CreateNewUser(s) // 创建新用户
	bituser, _ := proto.Marshal(mUser.GetInfo())
	s.Send((uint32)(PB.MSG_ID_USER_INFO), bituser) // 发送用户详细信息

	bitalluser, _ := proto.Marshal(g.GetAllUserData(mUser.Id))
	s.Send((uint32)(PB.MSG_ID_ALLPLAYERDATA), bitalluser) // 发送其它玩家信息

	bit_cmdcache, _ := proto.Marshal(g.GetAllUserLastCmdCache(mUser.Id))
	s.Send((uint32)(PB.MSG_ID_AllLASTCMDCACHE), bit_cmdcache) // 发送其它玩家操作缓存

	g.PlayerOnline(mUser) // 玩家上线

}

// 移动到
func (g *Game) C2S_MOVETO(s *MServer.MSession, bitmsg []byte) {
	m := PB.MoveTo{}
	proto.Unmarshal(bitmsg, &m)
	//fmt.Println("move", s.ID, m.Frompos.X, m.Frompos.Y)
	user := g.GetUser(s.ID)
	if user != nil {
		user.X = m.Frompos.X
		user.Y = m.Frompos.Y
		user.A = m.Frompos.Z
		user.MoveState = 1
		g.BroadCast((uint32)(PB.MSG_ID_MOVE_TO), bitmsg)
	}
}

//停止移动
func (g *Game) C2S_MOVESTOP(s *MServer.MSession, bitmsg []byte) {
	m := PB.MoveStop{}
	proto.Unmarshal(bitmsg, &m)

	//fmt.Println("stop", s.ID, m.Pos.X, m.Pos.Y)
	user := g.GetUser(s.ID)
	if user != nil {
		user.X = m.Pos.X
		user.Y = m.Pos.Y
		user.A = m.Pos.Z
		user.MoveState = 0
		g.BroadCast((uint32)(PB.MSG_ID_MOVE_STOP), bitmsg)
	}
}

// 系统信息PB
var sysinfo = PB.SysInfo{}

// 获得系统信息
func (g *Game) GetSysInfo() *PB.SysInfo {
	sysinfo.ServerId = g.server_id
	sysinfo.ServerRuntime = g.GetRunTime()
	sysinfo.CurFrame = g.GetTime()
	sysinfo.PlayerCount = g.user_count
	return &sysinfo
}
