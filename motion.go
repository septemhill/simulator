package simulator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
)

type Motion struct {
	motionMap map[string]int
}

type MotionManager struct {
	motion        *Motion
	self          *Self
	motionCh      chan *Motion
	modCmdCh      chan ModuleCommand
	modCmdHndlMap ModuleCommandHandlerMap
	//TODO: overMotionAction
}

var (
	initMotionCount = 0
)

var defaultMotionMap = Motion{
	motionMap: make(map[string]int),
}

func loadMotionList(filename string) ([]string, error) {
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

func (mm *MotionManager) Start() {
	go mm.taskHandler()
}

func (mm *MotionManager) Send(v interface{}) {
	switch val := v.(type) {
	case *Motion:
		mm.motionCh <- val
	case ModuleCommand:
		mm.modCmdCh <- val
	}
}

func (mm *MotionManager) updateMotion(m *Motion) {

}

func (mm *MotionManager) moduleCommandHandler(cmd ModuleCommand) {
	mm.modCmdHndlMap[cmd]()
}

func (mm *MotionManager) taskHandler() {
	for {
		select {
		case motion := <-mm.motionCh:
			mm.updateMotion(motion)
			//fmt.Println(motion)
		case cmd := <-mm.modCmdCh:
			mm.moduleCommandHandler(cmd)
			fmt.Println("Motion", cmd)
		}
	}
}

func NewMotionManager(s *Self) *MotionManager {
	cmdHndlMap := make(ModuleCommandHandlerMap, 0)

	cmdHndlMap[MODCMD_START] = func() { fmt.Println("motion start") }
	cmdHndlMap[MODCMD_SAVE] = func() { fmt.Println("motion save") }

	mm := &MotionManager{
		motionCh:      make(chan *Motion, 10),
		modCmdCh:      make(chan ModuleCommand, 5),
		modCmdHndlMap: cmdHndlMap,
	}

	motions, _ := loadMotionList(os.Getenv("GOPATH") + "/src/github.com/septemhill/humansim/simulator/motion.json")

	for _, motion := range motions {
		defaultMotionMap.motionMap[motion] = 0
	}

	return mm
}
