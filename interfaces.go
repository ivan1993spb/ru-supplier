package main

type InterfaceStarter func(ServIn, ConfIn, FiltIn) error

type ServIn interface {
	Start() error
	ShutDown() error
	IsRunning() bool
	RemoveCache() error
}

type ConfIn interface {
	SetFilterEnabled(bool)
	// SetHost(string)
	// SetPort(string)
	Save() error
}

type FiltIn interface {
	// SetExpsAll(PatternSet)
	// SetExpsOrderName(PatternSet)
	// SetExpsOKDP(PatternSet)
	// SetExpsOKPD(PatternSet)
	// SetExpsOrganisationName(PatternSet)
	// Save() error
}
