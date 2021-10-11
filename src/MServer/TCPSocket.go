// MServer 服务器 - TCP Socket
package MServer

import (
	"encoding/binary"
	"fmt"
	"net"
)

// TCPSocket开始 :5555
func (s *MServer) StartTCPSocket(addr string) {

	go func() {
		tcpSocktListener, _ := net.Listen("tcp", addr)

		defer tcpSocktListener.Close()

		for !s.bServerEnd {
			conn, err := tcpSocktListener.Accept()
			if err != nil {
				fmt.Println("\n", err.Error(), "use of closed network connection")
				continue
			}

			s.accept_TCPSocket_Session(conn)
		}
	}()

}

//接收session
func (s *MServer) accept_TCPSocket_Session(_conn net.Conn) {

	session := MSession{
		ConnType:       TCPSOCKET,
		conn:           _conn,
		ADDR:           _conn.RemoteAddr().String(),
		callback_error: s.session_error,
		callback_recv:  s.CB_recv,
	}

	s.NewSession(&session)

	go session.tcpSocket_recv()

}

//读取
// ||---固定消息头(2字节)---||--消息ID(4字节)--||--消息长度(4字节)(消息体去掉头尾)--||--消息体(N字节)--||--结束字符(2字节)--||
func (s *MSession) tcpSocket_recv() {

	buff := make([]byte, 0) //接收缓存, 用来储存截断的消息
	for {

		buffer := make([]byte, RECV_BUFF_LEN)
		n, err := s.conn.Read(buffer)

		if err != nil {
			if s.callback_error != nil {
				s.callback_error(s, err)
			} else {
				fmt.Println(s.conn.RemoteAddr().String(), " read error: ", err)
				s.conn.Close()
			}
			return
		}

		buff = s.BytesCombine(buff, buffer[:n])
		bufflen := len(buff)

		ihead := 0
		iend := 0
		for ihead >= 0 {
			ihead = indexOf(buff, ByteHead)

			if ihead < 0 || bufflen < ihead+10 {
				break
			}
			msglenBit := buff[ihead+4+2 : ihead+4+4+2]
			msglen := (int)(binary.LittleEndian.Uint32(msglenBit))
			iend = ihead + msglen + 2

			if iend+2 > bufflen {
				break
			}

			if buff[iend] != ByteEnd[0] || buff[iend+1] != ByteEnd[1] {
				break
			}

			if ihead >= 0 && iend > ihead {
				s.recvMessage(buff[ihead+2 : iend])
			}

			buff = buff[iend+2:]

		}
	}

}
