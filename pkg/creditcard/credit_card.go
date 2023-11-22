package creditcard

import (
	"time"
)

type CreditCard struct {
	Name      string `json:"name"`
	DaysToPay uint8  `json:"days-to-pay"`
	CutOffDay uint8  `json:"cutoff-day"`
}

func NewCreditCardFromStatementInfo(name string, lastStatementDate time.Time, lastStatementPaymentDate time.Time) CreditCard {
	diff := lastStatementPaymentDate.Sub(lastStatementDate)
	daysToPay := uint8(diff.Hours() / 24)

	return CreditCard{
		Name:      name,
		DaysToPay: daysToPay,
		CutOffDay: uint8(lastStatementDate.Day()),
	}

}
