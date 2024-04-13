package basic

type StateID int

const (
	EXIT      StateID = iota
	INIT      StateID = iota
	ATTEMPTER StateID = iota
	LEADER    StateID = iota
	FAILOVER  StateID = iota
	STOPPING  StateID = iota
)
