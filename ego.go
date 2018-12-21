package simulator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
)

type Ego struct {
	egoMap map[string]int
}

type EgoManager struct {
	ego           *Ego
	self          *Self
	eventCh       chan *Event
	modCmdCh      chan ModuleCommand
	modCmdHndlMap ModuleCommandHandlerMap
}

var (
	initEgoCount int32 = 0
)

var defaultEgoMap = Ego{
	egoMap: make(map[string]int),
}

func loadEgoList(filename string) ([]string, error) {
	fmt.Println(filename)
	var attrs []string
	fd, err := os.OpenFile(filename, syscall.O_RDONLY, 0664)

	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(fd)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &attrs)

	if err != nil {
		return nil, err
	}

	return attrs, nil
}

func (em *EgoManager) Start() {
	go em.taskHandler()
}

func (em *EgoManager) Send(v interface{}) {
	switch val := v.(type) {
	case *Event:
		em.eventCh <- val
	case ModuleCommand:
		em.modCmdCh <- val
	}
}

func (em *EgoManager) moduleCommandHandler(cmd ModuleCommand) {
	em.modCmdHndlMap[cmd]()
}

func (em *EgoManager) taskHandler() {
	for {
		select {
		case event := <-em.eventCh:
			em.self.modules[MODULE_MOTIONMGR].Send(event)
			//fmt.Println(event)
		case cmd := <-em.modCmdCh:
			em.moduleCommandHandler(cmd)
			fmt.Println("Ego", cmd)
		}
	}
}

func NewEgoManager(s *Self) *EgoManager {
	cmdHndlMap := make(ModuleCommandHandlerMap, 0)

	cmdHndlMap[MODCMD_START] = func() { fmt.Println("ego start") }
	cmdHndlMap[MODCMD_SAVE] = func() { fmt.Println("ego save") }

	em := &EgoManager{
		self:          s,
		eventCh:       make(chan *Event, 10),
		modCmdCh:      make(chan ModuleCommand, 5),
		modCmdHndlMap: cmdHndlMap,
	}

	egos, _ := loadEgoList(os.Getenv("GOPATH") + "/src/github.com/septemhill/humansim/simulator/ego.json")

	for _, ego := range egos {
		defaultEgoMap.egoMap[ego] = 0
	}

	return em
}
