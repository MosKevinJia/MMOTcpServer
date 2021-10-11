// MServer 服务器 - WebSocket
package MServer

import (
	//"bytes"
	"encoding/binary"
	"fmt"
	"net/http"
	"time"

	//"sync"

	"github.com/gorilla/websocket"
)

const (
	// 允许等待的写入时间
	writeWait = 0 //30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 0 //30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	//pingPeriod = (pongWait * 9) / 10

	ws_messageType = 2
)

var defWebsocketURL = "/ws"

var upgraderWebSocket = websocket.Upgrader{
	ReadBufferSize:  RECV_BUFF_LEN,
	WriteBufferSize: WRITE_BUFF_LEN,
	CheckOrigin: func(r *http.Request) bool { // 允许所有的CORS 跨域请求，正式环境可以关闭
		return true
	},
}

// WebSockt
func (m *MServer) StartTCPWebSocket(addrPort string) {
	go func() {
		http.HandleFunc(defWebsocketURL, m.accept_TCPWebSocket_Session)
		http.ListenAndServe(addrPort, nil)
	}()
}

//接收session
func (m *MServer) accept_TCPWebSocket_Session(resp http.ResponseWriter, req *http.Request) {

	fmt.Println("websocket accept")

	wsSocket, err := upgraderWebSocket.Upgrade(resp, req, nil)
	if err != nil {
		fmt.Println("升级为websocket失败", err.Error())
		return
	}

	session := MSession{
		ConnType:       TCPWEBSOCKET,
		wsconn:         wsSocket,
		ADDR:           wsSocket.RemoteAddr().String(),
		callback_error: m.session_error,
		callback_recv:  m.CB_recv,
	}

	m.NewSession(&session)

	go session.tcpWebSocket_recv()

}

//读取
// ||---固定消息头(2字节)---||--消息ID(4字节)--||--消息长度(4字节)--||--消息体(N字节)--||--结束字符(2字节\r\n)--||
func (s *MSession) tcpWebSocket_recv() {

	// 设置消息的最大长度
	s.wsconn.SetReadLimit(RECV_BUFF_LEN)
	//s.wsconn.SetReadDeadline(time.Now().Add(pongWait))
	//s.wsconn.SetWriteDeadline(time.Now().Add(writeWait))
	s.wsconn.SetReadDeadline(time.Time{})
	s.wsconn.SetWriteDeadline(time.Time{})

	for {

		_, buff, err := s.wsconn.ReadMessage()
		if err != nil {
			//websocket.IsUnexpectedCloseError(err, s.wsconn.CloseGoingAway, websocket.CloseAbnormalClosure)
			//fmt.Println("WebSocket消息读取出现错误", err.Error())
			//s.wsconn.Close()
			s.callback_error(s, err)
			return
		}
		//s.ws_messageType = mt

		bufflen := len(buff)
		ihead := indexOf(buff, ByteHead)
		//iend := indexOf(buff, ByteEnd)

		if ihead < 0 || bufflen < ihead+10 {
			break
		}

		msglenBit := buff[ihead+4+2 : ihead+4+4+2]
		msglen := (int)(binary.LittleEndian.Uint32(msglenBit))
		iend := ihead + msglen + 2

		if iend+2 > bufflen {
			break
		}

		//fmt.Println(buff[iend], ByteEnd[0], buff[iend] == ByteEnd[0])
		if buff[iend] != ByteEnd[0] || buff[iend+1] != ByteEnd[1] {
			break
		}

		if ihead >= 0 && iend > ihead {
			//fmt.Println(buff[ihead+2 : iend])
			s.recvMessage(buff[ihead+2 : iend])
		}

	}

}
