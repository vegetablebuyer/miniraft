package raft


type StateMachine interface {
	Save() ([]byte, error)
	Recover([]byte) error
}
