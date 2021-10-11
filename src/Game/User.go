// 用户
package Game

import (
	"TestMMO2/src/Game/PB"
	"TestMMO2/src/MServer"
)

// User
type User struct {
	Id   int32
	Name string
	Pwd  string
	Type int32

	HP     int32
	MP     int32
	HP_MAX int32
	MP_MAX int32

	// 位置
	X float32
	Y float32
	// 方向
	A float32

	MoveState int32 //是否移动 0=停 1坐标移动  2方向移动
	Speed     float32
	FRAME int32

	Session *MServer.MSession

	Info PB.UserInfo
}

// 获得用户基础信息
func (u *User) UpdateInfo() {
	u.Info.Id = u.Id
	u.Info.Name = u.Name
	u.Info.Type = u.Type
	u.Info.Pos = &PB.Vec3{X: u.X, Y: u.Y, Z: u.A}
	u.Info.MoveState = u.MoveState
	u.Info.Hp = u.HP
	u.Info.Mp = u.MP
	u.Info.HpMax = u.Info.HpMax
	u.Info.MpMax = u.Info.MpMax
}

// 获得用户基础信息
func (u *User) GetInfo() *PB.UserInfo {
	u.UpdateInfo()
	return &u.Info
}
