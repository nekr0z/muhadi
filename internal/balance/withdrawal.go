package balance

import (
	"errors"
	"time"
)

type Withdrawal struct {
	UserName string
	OrderID  int
	Amount   float64
	At       time.Time
}

func NewWithdrawal(userName string, order int, amount float64) *Withdrawal {
	return &Withdrawal{
		UserName: userName,
		OrderID:  order,
		Amount:   amount,
		At:       time.Now(),
	}
}

var ErrNotEnoughFunds = errors.New("not enough funds")
