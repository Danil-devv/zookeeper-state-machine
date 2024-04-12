package number

type State int

const (
	EXIT      State = iota
	INIT      State = iota
	ATTEMPTER State = iota
	LEADER    State = iota
	FAILOVER  State = iota
	STOPPING  State = iota
)
