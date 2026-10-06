package ledger

type Direction string

const (
	DirectionDebit  Direction = "DEBIT"
	DirectionCredit Direction = "CREDIT"
)

func (d Direction) IsValid() bool {
	return d == DirectionDebit || d == DirectionCredit
}
