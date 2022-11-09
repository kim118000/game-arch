package network

import "fmt"

type ClusterAddr struct {
	Name      string
	Addr      string
	Id        uint32
	Reconnect bool
	IsHand    bool
	GroupId   int
}

func (c *ClusterAddr) String() string {
	return fmt.Sprintf("[name=%s,addr=%s,id=%d]", c.Name, c.Addr, c.Id)
}