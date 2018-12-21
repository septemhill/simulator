package simulator

type Observer interface {
	Update(interface{})
}

type Subject interface {
	Attach(Observer)
	Detach(Observer)
	Notify()
}
