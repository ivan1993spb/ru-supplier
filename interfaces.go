package main

type InterfaceStarter func(ServIn, ConfIn, FiltIn, CachIn) error

type ServIn interface {
	Start() error
	ShutDown() error
	IsRunning() bool
}

type ConfIn interface {
	SetFilterEnabled(bool)
	SetHost(string)
	SetPort(string)
	Save() error
}

type FiltIn interface {
	SetExpsAll(PatternSet)
	SetExpsOrderName(PatternSet)
	SetExpsOKDP(PatternSet)
	SetExpsOKPD(PatternSet)
	SetExpsOrganisationName(PatternSet)
	Save() error
}

type CachIn interface {
	Save() error
	Remove() error
}
