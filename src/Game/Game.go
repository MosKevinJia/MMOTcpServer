// 游戏主程序
package Game

import (
	"TestMMO2/src/MServer"
	"fmt"
	"strconv"

	"sync"
	"time"
)

type Game struct {
	server_id          int32 // ServerID
	user_maxcount      int32 // 最大游戏人数
	user_count         int32 // 当前游戏人数
	user_defid         int32 // 玩家起始ID
	start_run_time     int64 // 开始运行时间
	net_mode_socket    bool  // 是否使用普通Socket模式
	net_mode_websocket bool  // 是否使用WebSocket模式

	UserDict sync.Map // 用户

	RouterDict map[uint32]CallBackRouter // 路由字典

	server *MServer.MServer

	mux sync.RWMutex
}

func init() {

}

// 新创建
func New() *Game {
	game := Game{}
	game.user_maxcount = 1000 // 最大连接数
	game.user_defid = 0       // 起始ID
	game.InitRouter()

	// MServer
	game.server = MServer.New()
	game.server.CB_recv = game.SessionRecv
	game.server.CB_newclient = game.SessionAccpet
	game.server.CB_error = game.SessionError

	return &game
}

// 设置Id
func (g *Game) SetID(id int32) {
	g.server_id = id
}

// Listen TCPSocket
func (g *Game) StartTCPSocket(port int) {
	g.net_mode_socket = true
	g.server.StartTCPSocket(":" + strconv.Itoa(port))
	fmt.Println("Socket Listen", port)

}

// Listen WebSockt
func (g *Game) StartTCPWebSocket(port int) {
	g.net_mode_websocket = true
	g.server.StartTCPWebSocket(":" + strconv.Itoa(port))
	fmt.Println("WebSockt Listen", port)
}

// Game Start
func (g *Game) Start() {
	fmt.Println("Game Server Start")
	g.start_run_time = g.GetTime()
}

// 设置同时最大人数
func (g *Game) SetMaxUserLen(maxlen int32) {
	g.user_maxcount = maxlen
}

// 关闭
func (g *Game) Close() {
	fmt.Println("Game Server Close")
}

// 获取帧
func (g Game) GetTime() int64 {
	return time.Now().UnixNano()
}

// 获取运行时间
func (g Game) GetRunTime() int64 {
	return g.GetTime() - g.start_run_time
}
