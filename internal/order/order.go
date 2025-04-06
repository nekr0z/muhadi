package order

import (
	"errors"
	"time"
)

type Status int

const (
	StatusNew Status = iota
	StatusProcessing
	StatusInvalid
	StatusDone
)

func (s Status) String() string {
	switch s {
	case StatusNew:
		return "NEW"
	case StatusProcessing:
		return "PROCESSING"
	case StatusInvalid:
		return "INVALID"
	case StatusDone:
		return "PROCESSED"
	default:
		return "UNKNOWN"
	}
}

type Order struct {
	ID         int
	UserName   string
	Status     Status
	Accrual    float64
	LastUpdate time.Time
	CreatedAt  time.Time
}

func New(id int, userName string) *Order {
	return &Order{
		ID:         id,
		UserName:   userName,
		Status:     StatusNew,
		Accrual:    0,
		CreatedAt:  time.Now(),
		LastUpdate: time.Now(),
	}
}

func (o *Order) UpdateStatus(status Status) {
	o.Status = status
	o.LastUpdate = time.Now()
}

func (o *Order) UpdateAccrual(accrual float64) {
	o.Accrual = accrual
	o.LastUpdate = time.Now()
}

var ErrOrderNumberInvalid = errors.New("order number invalid")
