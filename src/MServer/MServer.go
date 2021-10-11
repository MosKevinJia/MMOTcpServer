// MServer 服务器
package MServer

import (
	"bytes"
	"encoding/binary"

	//"fmt"
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

type NET_TYPE int32 // 连接方式

const (
	TCPSOCKET      NET_TYPE = 0    // Tcp Socket
	TCPWEBSOCKET   NET_TYPE = 1    // Tcp WebSocket
	RECV_BUFF_LEN           = 1024 // 接收缓冲区大小
	WRITE_BUFF_LEN          = 1024
)

var ByteHead = []byte("ms")
var ByteEnd = []byte("!&")

type MServer struct {
	mAllSessions         sync.Map
	bServerEnd           bool
	tcpSocketListener    net.Listener    // tcpsockt监听
	tcpWebSocketListener *websocket.Conn // websocket监听

	CB_newclient CallBack_NewClient // 新添加
	CB_recv      CallBack_Recv      // 接收
	CB_error     CallBack_Error     //
}

type CallBack_NewClient func(s *MSession)
type CallBack_Recv func(s *MSession, msgid uint32, bit_msg []byte)
type CallBack_Error func(s *MSession, err error)

func init() {
	//fmt.Println("head", ByteHead)
	//fmt.Println("end", ByteEnd)
}

// 新创建
func New() *MServer {
	s := &MServer{}
	return s
}

// 广播
func (s *MServer) BroadCast(msgid uint32, bitmsg []byte) {

	var bid = make([]byte, 4)
	binary.LittleEndian.PutUint32(bid, msgid) // 加入Msg ID
	var blen = make([]byte, 4)
	binary.LittleEndian.PutUint32(blen, (uint32)(len(bitmsg)+8)) // 加入长度Len

	bitm := BytesCombine(ByteHead, bid, blen, bitmsg, ByteEnd) // 合并

	s.mAllSessions.Range(func(k, v interface{}) bool {
		v.(*MSession).SendByte(bitm)
		return true
	})

}

// 设置head
func (s *MServer) SetHead(head []byte) {
	if len(head) == 4 {
		ByteHead = head
	}
}

// 设置End
func (s *MServer) SetEnd(end []byte) {
	if len(end) == 4 {
		ByteEnd = end
	}
}

// 新Session
func (s *MServer) NewSession(session *MSession) {
	if s.CB_newclient != nil {
		s.CB_newclient(session)
	}

	s.mAllSessions.Store(session.ADDR, session)

}

// 服务端主动关闭Session
func (s *MServer) CloseSession(session *MSession) {
	if session != nil {
		if session.ConnType == TCPSOCKET {
			session.conn.Close()
		} else if session.ConnType == TCPWEBSOCKET {
			session.wsconn.Close()
		}

		s.mAllSessions.Delete(session.ADDR)
	}
}

// session 错误(断线或客户端主动关闭)
func (s *MServer) session_error(session *MSession, err error) {

	if s.CB_error != nil {
		s.CB_error(session, err)
	}

	if session.ConnType == TCPSOCKET {
		session.conn.Close()
	} else if session.ConnType == TCPWEBSOCKET {
		session.wsconn.Close()
	}

	s.mAllSessions.Delete(session.ADDR)

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
