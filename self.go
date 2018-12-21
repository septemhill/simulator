package simulator

import (
	"time"
)

type Self struct {
	Human
	modules map[string]Module
}

func (s *Self) Noop() {
	tick := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-tick.C:
			for _, module := range s.modules {
				module.Send(MODCMD_SAVE)
			}
		}
	}
}

func NewSelf() *Self {
	s := &Self{
		modules: make(map[string]Module),
	}

	s.modules[MODULE_MOTIONMGR] = NewEmotionManager(s)
	s.modules[MODULE_EGOMDR] = NewEgoManager(s)
	s.modules[MODULE_COMMUMGR] = NewCommunicateManager(s)

	for _, module := range s.modules {
		module.Start()
	}

	go s.Noop()

	return s
}

//func (s *Self) saveConnectedHuman() {
//	for _, other := range s.connected {
//		if reflect.DeepEqual(*other, Other{}) {
//		}
//	}
//}
//
