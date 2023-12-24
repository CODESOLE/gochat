package core

import "net"

type Payload struct {
  Msg []byte
  Conn net.Conn
  IpAddr string
  SendTime string
}
