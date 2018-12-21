package simulator

type ModuleCommand int32

type Module interface {
	Start()
	Send(interface{})
	//TaskConsumer()
}

type ModuleCommandHandler func()

type ModuleCommandHandlerMap map[ModuleCommand]ModuleCommandHandler

const (
	MODULE_MOTIONMGR = "mmMgr"
	MODULE_EGOMDR    = "egoMdr"
	MODULE_COMMUMGR  = "commuiMgr"
)

const (
	MODCMD_START ModuleCommand = iota
	MODCMD_SAVE
	MODCMD_MAXIMUM
)
