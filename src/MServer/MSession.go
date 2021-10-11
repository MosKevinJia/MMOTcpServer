package MServer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

type MSession struct {
	ConnType       NET_TYPE // 连接类型
	ID             int32
	ADDR           string
	callback_error CallBack_Error
	callback_recv  CallBack_Recv
	conn           net.Conn
	wsconn         *websocket.Conn
	mux            sync.RWMutex
	cmdCache       *[]byte // 命令缓存
	cmdCacheId     uint32
}

// 验证并处理消息
func (s *MSession) recvMessage(bit_msg []byte) {

	blen := (uint32)(len(bit_msg))

	bitid := bit_msg[:4]
	bitlen := bit_msg[4:8]

	id := binary.LittleEndian.Uint32(bitid)
	le := binary.LittleEndian.Uint32(bitlen)

	if le > blen {
		fmt.Println("Recv 数据出错: 半包 ", blen, id, le)
		fmt.Println(bit_msg)
		return
	}

	bmsg := bit_msg[8:]
	s.cmdCache = &bmsg
	s.cmdCacheId = id

	if s.callback_recv != nil {
		s.callback_recv(s, id, bmsg)
	}

}

/// <summary>
/// 查找Byte字符
/// </summary>
/// <param name="msg"></param>
/// <param name="key"></param>
/// <returns></returns>
func indexOf(b []byte, bb []byte) int {
	ib := len(b)
	ibb := len(bb)
	if b == nil || bb == nil || ib < 1 || ibb < 1 || ib < ibb {
		return -1
	}

	i := 0
	j := 0
	for i = 0; i < ib-ibb+1; i++ {
		if b[i] == bb[0] {
			for j = 1; j < ibb; j++ {
				if b[i+j] != bb[j] {
					break
				}
			}
			if j == ibb {
				return i
			}
		}
	}

	return -1
}

// 直接发送
func (s *MSession) SendByte(bitmsg []byte) {

	s.mux.Lock() //加锁
	if s.ConnType == TCPSOCKET {
		_, err := s.conn.Write(bitmsg)
		if err != nil {
			fmt.Println(err)
		}
	} else if s.ConnType == TCPWEBSOCKET {
		s.wsconn.WriteMessage(ws_messageType, bitmsg)
	}
	s.mux.Unlock() //解锁
}

// 使用套接字发送
func (s *MSession) Send(msgid uint32, bitmsg []byte) {

	var bid = make([]byte, 4)
	binary.LittleEndian.PutUint32(bid, msgid)
	var blen = make([]byte, 4)
	binary.LittleEndian.PutUint32(blen, (uint32)(len(bitmsg)+8))
	bitm := s.BytesCombine(ByteHead, bid, blen, bitmsg, ByteEnd)
	s.SendByte(bitm)

}

// 获得命令缓存
func (s *MSession) GetLastCmdCache() (uint32, *[]byte) {
	return s.cmdCacheId, s.cmdCache
}

//BytesCombine 多个[]byte数组合并成一个[]byte
func (session *MSession) BytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}
