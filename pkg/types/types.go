package types

type Money int

type PaymentCategory string

type PaymentStatus string

type Phone string

const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "Fail"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

type Favorite struct {
	ID        string
	Name      string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
}

func (p *Payment) GetPayment() Payment {
	return Payment{
		Amount:    p.Amount,
		AccountID: p.AccountID,
		Category:  p.Category,
	}
}
