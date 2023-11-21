package creditcard

import (
  "time"
  "encoding/json"
)

type CreditCardDTO struct {
  Name string `json:"name"`
  DaysToPay json.Number `json:"days-to-pay"`
  CutOffDay json.Number `json:"cutoff-day"`
}


type CreditCard struct {
  Name string 
  DaysToPay uint8
  CutOffDay uint8
}

func NewCreditCardFromStatementInfo(name string, lastStatementDate time.Time, lastStatementPaymentDate time.Time) CreditCard {
  diff := lastStatementPaymentDate.Sub(lastStatementDate)
  daysToPay := uint8(diff.Hours() / 24)


  return CreditCard {
    Name: name,
    DaysToPay: daysToPay,
    CutOffDay: uint8(lastStatementDate.Day()),
  }

}
