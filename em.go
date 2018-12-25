package simulator

import (
	"fmt"
	"math/rand"

	"github.com/septemhill/fion"
)

type EmotionValue int32
type EmotionThreshold int32
type EmotionAction func()

type Emotion struct {
	Name           []byte
	Nickname       string
	Difference     chan EmotionValue
	Value          EmotionValue
	ThresValueMap  []EmotionValue
	ThresActionMap map[EmotionThreshold][]EmotionAction
}

const (
	EMTHRES_NO_EMOTION EmotionThreshold = iota
	EMTHRES_UNDER_EMOTION
	EMTHRES_ON_EMOTION
	EMTHRES_OVER_EMOTION
	EMTHRES_MAX_EMOTION
)

func (e *Emotion) Start() {
	go e.emotionTaskHandler()
}

func (e *Emotion) Send(v interface{}) {
	switch v := v.(type) {
	case EmotionValue:
		e.Difference <- v
	default:
		fion.Warn("unknown type")
	}
}

func (e *Emotion) checkEmotionThreshold() EmotionThreshold {
	thres := EMTHRES_NO_EMOTION

	for i := 0; i < len(e.ThresValueMap); i++ {
		if e.Value < e.ThresValueMap[i] {
			break
		}
		thres = EmotionThreshold(i)
	}

	return thres
}

func (e *Emotion) doEmotionAction(thres EmotionThreshold) {
	actions := e.ThresActionMap[thres]

	if len(actions) == 0 {
		return
	}

	rand.Shuffle(len(actions), func(i, j int) {
		actions[i], actions[j] = actions[j], actions[j]
	})

	actions[0]()
}

func (e *Emotion) emotionTaskHandler() {
	for {
		select {
		case val := <-e.Difference:
			e.Value += val
			fmt.Println(e.Value)
			thres := e.checkEmotionThreshold()
			e.doEmotionAction(thres)
		}
	}
}

func NewEmotion(name []byte, nickname string) *Emotion {
	e := &Emotion{
		Name:           name,
		Nickname:       nickname,
		Difference:     make(chan EmotionValue, 10),
		ThresValueMap:  make([]EmotionValue, EMTHRES_MAX_EMOTION),
		ThresActionMap: make(map[EmotionThreshold][]EmotionAction),
	}

	e.ThresValueMap[EMTHRES_UNDER_EMOTION] = 20
	e.ThresValueMap[EMTHRES_ON_EMOTION] = 50
	e.ThresValueMap[EMTHRES_OVER_EMOTION] = 80

	for i := 0; i < int(EMTHRES_MAX_EMOTION); i++ {
		e.ThresActionMap[EmotionThreshold(i)] = make([]EmotionAction, 0)
	}

	e.ThresActionMap[EMTHRES_UNDER_EMOTION] = append(e.ThresActionMap[EMTHRES_UNDER_EMOTION], func() {
		fmt.Println("Hi, Septem. Do under emtion action")
	})

	return e
}
