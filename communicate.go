package simulator

import (
	"fmt"
	"net"

	"github.com/septemhill/fion"
)

//type CommuType int
//
//const (
//	COMMU_QUESTION CommuType = iota
//	COMMU_ANSWER
//	COMMU_PERSONINTRO
//	COMMU_NETWORKIDASK
//	COMMU_REALIDASK
//)
//
//type Message struct {
//	Type    CommuType
//	Content interface{}
//}
//
//func Talk(w io.Writer, msg []byte) {
//	w.Write(msg)
//}
//

type Other struct {
	Human
	answered int
}

type CommunicateManager struct {
	self          *Self
	connCh        chan net.Conn
	modCmdCh      chan ModuleCommand
	listener      net.Listener
	connected     []*Other
	modCmdHndlMap ModuleCommandHandlerMap
}

const (
	connectionMaxLimitation = 1000
)

func (c *CommunicateManager) Start() {
	go c.enableCommunication(c.listener)
	go c.taskHandler()
}

func (c *CommunicateManager) Send(v interface{}) {
	switch val := v.(type) {
	case net.Conn:
		c.connCh <- val
	case ModuleCommand:
		c.modCmdCh <- val
	}
}

func (c *CommunicateManager) newConnHandler(conn net.Conn) {
	if len(c.connected) >= connectionMaxLimitation {
		conn.Close()
		return
	}

	o := &Other{}
	c.connected = append(c.connected, o)
}

func (c *CommunicateManager) moduleCommandHandler(cmd ModuleCommand) {
	c.modCmdHndlMap[cmd]()
}

func (c *CommunicateManager) taskHandler() {
	for {
		select {
		case conn := <-c.connCh:
			c.newConnHandler(conn)
			//fmt.Println(conn)
		case cmd := <-c.modCmdCh:
			c.moduleCommandHandler(cmd)
			fmt.Println("Commu", cmd)
		}
	}
}

func (c *CommunicateManager) enableCommunication(listener net.Listener) {
	for {
		conn, err := listener.Accept()

		if err != nil {
			fion.Warn(err)
			continue
		}

		c.connCh <- conn
	}
}

func newCommnuicateListener() (net.Listener, error) {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		return nil, err
	}

	return listener, nil
}

func NewCommunicateManager(s *Self) *CommunicateManager {
	listener, err := newCommnuicateListener()
	cmdHndlMap := make(ModuleCommandHandlerMap, 0)

	cmdHndlMap[MODCMD_START] = func() { fmt.Println("commu start") }
	cmdHndlMap[MODCMD_SAVE] = func() { fmt.Println("commu save") }

	if err != nil {
		panic("failed to create communicate listener")
	}

	comm := &CommunicateManager{
		self:          s,
		connCh:        make(chan net.Conn, 10),
		modCmdCh:      make(chan ModuleCommand, 5),
		connected:     make([]*Other, 0),
		listener:      listener,
		modCmdHndlMap: cmdHndlMap,
	}

	return comm
}
