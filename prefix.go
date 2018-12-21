package simulator

type DBDataPrefix byte

const (
	DBPREFIX_REALID     DBDataPrefix = iota //Real Identity ID
	DBPREFIX_QNA                            //Question and Answer
	DBPREFIX_PERSONINFO                     //Other Human Infomation
	DEPREFIX_MONTIONCHG                     //Motion Change Record
	DBPREFIX_EGO                            //Latest Ego State
	DBPREFIX_EGOLIST                        //Existed Ego List
	DBPREFIX_MOTIONLIST                     //Existed Motion List
)
