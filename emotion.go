package simulator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
)

const (
	EmotionUnderThreshold = 20
	EmotionOnThreshold    = 50
	EmotionOverThreshold  = 70
)

type EmotionalAction func()

type EmotionB struct {
	emotionMap map[string]int
}

type EmotionManager struct {
	emotion        *EmotionB
	self           *Self
	emotionCh      chan *EmotionB
	checkEmotionCh chan *EmotionB
	modCmdCh       chan ModuleCommand
	modCmdHndlMap  ModuleCommandHandlerMap
	//TODO: overEmotionAction
}

var (
	initEmotionCount = 0
)

var defaultEmotionMap = EmotionB{
	emotionMap: make(map[string]int),
}

func loadEmotionList(filename string) ([]string, error) {
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

func (mm *EmotionManager) Start() {
	go mm.taskHandler()
}

func (mm *EmotionManager) Send(v interface{}) {
	switch val := v.(type) {
	case *EmotionB:
		mm.emotionCh <- val
	case ModuleCommand:
		mm.modCmdCh <- val
	}
}

func (mm *EmotionManager) updateEmotion(m *EmotionB) {

}

func (mm *EmotionManager) moduleCommandHandler(cmd ModuleCommand) {
	mm.modCmdHndlMap[cmd]()
}

func (mm *EmotionManager) taskHandler() {
	for {
		select {
		case emotion := <-mm.emotionCh:
			mm.updateEmotion(emotion)
			//fmt.Println(emotion)
		case cmd := <-mm.modCmdCh:
			mm.moduleCommandHandler(cmd)
			fmt.Println("EmotionB", cmd)
		}
	}
}

func NewEmotionManager(s *Self) *EmotionManager {
	cmdHndlMap := make(ModuleCommandHandlerMap, 0)

	cmdHndlMap[MODCMD_START] = func() { fmt.Println("emotion start") }
	cmdHndlMap[MODCMD_SAVE] = func() { fmt.Println("emotion save") }

	mm := &EmotionManager{
		emotionCh:     make(chan *EmotionB, 10),
		modCmdCh:      make(chan ModuleCommand, 5),
		modCmdHndlMap: cmdHndlMap,
	}

	emotions, _ := loadEmotionList(os.Getenv("GOPATH") + "/src/github.com/septemhill/humansim/simulator/emotion.json")

	for _, emotion := range emotions {
		defaultEmotionMap.emotionMap[emotion] = 0
	}

	return mm
}
