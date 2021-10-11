// 用户
package Game

import (
	"TestMMO2/src/Game/PB"
	"TestMMO2/src/MServer"

	//"fmt"

	"math/rand"
	"strconv"

	proto "github.com/gogo/protobuf/proto"
)

// 创建新用户
func (g *Game) CreateNewUser(session *MServer.MSession) *User {

	var ID int32
	g.user_defid++
	ID = 1000 + g.user_defid // 用户ID
	session.ID = ID

	user := User{
		Id:      ID,
		Name:    "player" + strconv.FormatInt(int64(ID), 10),
		Type:    0,
		X:       rand.Float32()*10 - 5,
		Y:       rand.Float32()*10 - 5,
		A:       0,
		HP:      100,
		MP:      100,
		HP_MAX:  100,
		MP_MAX:  100,
		Speed:   30,
		Session: session,
		Info:    PB.UserInfo{},
	}

	user.UpdateInfo()

	return &user

}

// 从数据库加载用户
func (g *Game) LoadDBUser(id int32, session *MServer.MSession) *User {
	return nil
}

// 保存用户数据
func (g *Game) SaveDBUser(id int32, session *MServer.MSession) error {

	return nil
}

//
func (g *Game) GetUser(id int32) *User {
	u, ok := g.UserDict.Load(id)
	if ok {
		return u.(*User)
	}
	return nil
}

// 玩家上线
func (g *Game) PlayerOnline(user *User) {

	g.UserDict.Store(user.Id, user)

	g.mux.Lock()
	g.user_count++
	g.mux.Unlock()
	p := PB.PlayerOnline{
		Info:  user.GetInfo(),
		Frame: g.GetTime(),
	}
	bitmsg, _ := proto.Marshal(&p)

	g.BroadCast((uint32)(PB.MSG_ID_PLAYER_ONLINE), bitmsg)

}

// 玩家下线
func (g *Game) PlayerOffline(id int32) {

	_, ok := g.UserDict.LoadAndDelete(id)

	if !ok {
		return
	}

	g.mux.Lock()
	g.user_count--
	g.mux.Unlock()

	p := PB.PlayerOffline{
		Id:    id,
		Frame: g.GetTime(),
	}
	bitmsg, _ := proto.Marshal(&p)

	g.BroadCast((uint32)(PB.MSG_ID_PLAYER_OFFLINE), bitmsg)
}

// 获得所有Player UserData
func (g *Game) GetAllUserData(skip_id int32) *PB.AllPlayerData {

	arr := make([]*PB.UserInfo, 0)
	count := 0
	g.UserDict.Range(func(k, v interface{}) bool {
		if v.(*User).Id != skip_id {
			u := v.(*User)
			info := u.GetInfo()
			arr = append(arr, info)
			count++
		}
		return true
	})

	p := PB.AllPlayerData{
		Players: arr,
		Frame:   g.GetTime(),
	}

	return &p
}

// 获得所有Player UserData
func (g *Game) GetAllUserLastCmdCache(skip_id int32) *PB.AllLastCmdCache {

	arr := make([]*PB.LastCmdCache, 0)
	count := 0
	g.UserDict.Range(func(k, v interface{}) bool {
		if v.(*User).Id != skip_id {
			u := v.(*User)

			msgid, byteCmd := u.Session.GetLastCmdCache()

			if msgid != (uint32)(PB.MSG_ID_LOGIN) {
				//fmt.Println("byteCmd ", byteCmd)
				cache := PB.LastCmdCache{Id: (int32)(msgid), Cache: *byteCmd}
				arr = append(arr, &cache)
				count++
			}
		}
		return true
	})

	p := PB.AllLastCmdCache{
		Caches: arr,
	}

	return &p
}
