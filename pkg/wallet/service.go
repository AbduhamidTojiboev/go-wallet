package wallet

import (
	"errors"
	"github.com/AbduhamidTojiboev/go-wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("Phone already registered")
var ErrAccountNotFound = errors.New("Account not found")
var ErrAmountMustBePositive = errors.New("Amount must be greater than zero")
var ErrNotEnoughBalance = errors.New("Insufficient funds on balance")
var ErrPaymentNotFound = errors.New("Payment not found")
var ErrFavoriteNotFound = errors.New("Favorite not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Balance: 0,
		Phone:   phone,
	}

	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)

	if err != nil {
		return err
	}

	account.Balance += amount

	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	account, err := s.FindAccountByID(accountID)

	if err != nil {
		return nil, err
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}
	account.Balance -= amount
	paymentId := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentId,
		AccountID: accountID,
		Amount:    amount,
		Status:    types.PaymentStatusInProgress,
		Category:  category,
	}
	s.payments = append(s.payments, payment)

	return payment, nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	var payment *types.Payment
	for _, pay := range s.payments {
		if pay.ID == paymentID {
			payment = pay
			break
		}
	}

	if payment == nil {
		return nil, ErrPaymentNotFound
	}

	return payment, nil
}

func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	var favorite *types.Favorite
	for _, fav := range s.favorites {
		if fav.ID == favoriteID {
			favorite = fav
			break
		}
	}

	if favorite == nil {
		return nil, ErrFavoriteNotFound
	}

	return favorite, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	return account, nil
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return err
	}

	account, err := s.FindAccountByID(payment.AccountID)

	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount

	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return nil, err
	}

	return s.Pay(payment.AccountID, payment.Amount, payment.Category)
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return nil, err
	}

	favoriteId := uuid.New().String()
	favorite := &types.Favorite{
		ID:        favoriteId,
		Name:      name,
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, favorite)

	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)

	if err != nil {
		return nil, err
	}

	return s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
}
