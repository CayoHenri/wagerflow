package money

type Money struct {
	amount   int64
	currency Currency
}

func New(amount int64, currency Currency) Money {
	return Money{
		amount:   amount,
		currency: currency,
	}
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() Currency {
	return m.currency
}
