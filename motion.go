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

type Emotion struct {
	emotionMap map[string]int
}

type EmotionManager struct {
	emotion       *Emotion
	self          *Self
	emotionCh     chan *Emotion
	modCmdCh      chan ModuleCommand
	modCmdHndlMap ModuleCommandHandlerMap
	//TODO: overEmotionAction
}

var (
	initEmotionCount = 0
)

var defaultEmotionMap = Emotion{
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
	case *Emotion:
		mm.emotionCh <- val
	case ModuleCommand:
		mm.modCmdCh <- val
	}
}

func (mm *EmotionManager) updateEmotion(m *Emotion) {

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
			fmt.Println("Emotion", cmd)
		}
	}
}

func NewEmotionManager(s *Self) *EmotionManager {
	cmdHndlMap := make(ModuleCommandHandlerMap, 0)

	cmdHndlMap[MODCMD_START] = func() { fmt.Println("emotion start") }
	cmdHndlMap[MODCMD_SAVE] = func() { fmt.Println("emotion save") }

	mm := &EmotionManager{
		emotionCh:     make(chan *Emotion, 10),
		modCmdCh:      make(chan ModuleCommand, 5),
		modCmdHndlMap: cmdHndlMap,
	}

	emotions, _ := loadEmotionList(os.Getenv("GOPATH") + "/src/github.com/septemhill/humansim/simulator/emotion.json")

	for _, emotion := range emotions {
		defaultEmotionMap.emotionMap[emotion] = 0
	}

	return mm
}
