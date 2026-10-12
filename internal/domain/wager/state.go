package wager

type State string

const (
	StatePending          State = "PENDING"
	StatePendingReference State = "PENDING_REFERENCE"
	StateProcessed        State = "PROCESSED"
	StateRejected         State = "REJECTED"
	StateFailed           State = "FAILED"
)

func (s State) IsValid() bool {
	switch s {
	case StatePending, StatePendingReference, StateProcessed, StateRejected, StateFailed:
		return true
	default:
		return false
	}
}

func (s State) IsTerminal() bool {
	switch s {
	case StateProcessed, StateRejected, StateFailed:
		return true
	default:
		return false
	}
}
