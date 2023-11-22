package creditcard

import (
	"testing"
	"time"
)

func Test_NewCreditCardFromStatement(t *testing.T) {
	timeFormat := "2006-01-02 15:04:05"

	lastStatementDate, _ := time.Parse(timeFormat, "2023-10-11 00:00:00")
	lastStatementPayment, _ := time.Parse(timeFormat, "2023-11-06 00:00:00")

	cc := NewCreditCardFromStatementInfo("test", lastStatementDate, lastStatementPayment)

	const expDaysToPay uint8 = 26

	if cc.DaysToPay != expDaysToPay {
		t.Errorf("Incorrect number of days to pay, expected %d got %d", expDaysToPay, cc.DaysToPay)
	}

	const expCutOffDay uint8 = 11

	if cc.CutOffDay != expCutOffDay {
		t.Errorf("Incorrect cut off day, expected %d got %d", expCutOffDay, cc.CutOffDay)
	}
}
