package wallet

import (
	"fmt"
	"github.com/AbduhamidTojiboev/go-wallet/pkg/types"
	"reflect"
	"testing"
)

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccountWithBalance(phone types.Phone, amount types.Money) (*types.Account, error) {
	acc, err := s.RegisterAccount(phone)

	if err != nil {
		return nil, fmt.Errorf("can't register account, error %v", err)
	}

	err = s.Deposit(acc.ID, amount)

	if err != nil {
		return nil, fmt.Errorf("can't deposit account, error %v", err)
	}

	return acc, nil
}

func (s *testService) addDefaultAccountAndPayment() (*types.Account, *types.Payment, error) {
	phone := types.Phone("+992927894561")
	amountAccount := types.Money(100000)
	amountPay := types.Money(8)
	categoryPay := types.PaymentCategory("auto")
	account, err := s.addAccountWithBalance(phone, amountAccount)

	if err != nil {
		return nil, nil, err
	}

	payment, err := s.Pay(account.ID, amountPay, categoryPay)
	return account, payment, nil
}

func TestService_RegisterAccount_AlreadyRegistered(t *testing.T) {
	s := &Service{}
	_, err := s.RegisterAccount("+992927894561")

	if err != nil {
		t.Errorf("RegisterAccount(): can't register account, error %v", err)
		return
	}

	_, err = s.RegisterAccount("+992927894561")
	if err == nil {
		t.Errorf("RegisterAccount(): can't register account, error %v", err)
		return
	}
}

func TestService_RegisterAccount_Success(t *testing.T) {
	s := &Service{}
	phone := types.Phone("+992927894561")
	gotAccount, err := s.RegisterAccount(phone)

	if err != nil {
		t.Errorf("RegisterAccount(): can't register account, error %v", err)
		return
	}
	account := &types.Account{
		ID:      1,
		Phone:   phone,
		Balance: 0,
	}
	if !reflect.DeepEqual(gotAccount, account) {
		t.Errorf("RegisterAccount(): can't register account, error %v", err)
		return
	}
}

func TestService_Deposit_Success(t *testing.T) {
	s := &Service{}
	account, err := s.RegisterAccount("+992927894561")

	if err != nil {
		t.Errorf("RegisterAccount(): can't register account, error %v", err)
		return
	}

	err = s.Deposit(account.ID, 10)

	if err != nil {
		t.Errorf("Deposit(): can't deposit account, error %v", err)
		return
	}
}

func TestService_Deposit_AmountNegative(t *testing.T) {
	s := &Service{}
	account, err := s.RegisterAccount("+992927894561")

	if err != nil {
		t.Errorf("RegisterAccount(): can't register account, error %v", err)
		return
	}

	err = s.Deposit(account.ID, 0)

	if err == nil {
		t.Errorf("Deposit(): can't deposit account, error %v", err)
		return
	}
}

func TestService_Deposit_CannotAccount(t *testing.T) {
	s := &Service{}
	err := s.Deposit(2, 10)

	if err == nil {
		t.Errorf("Deposit(): can't deposit account, error %v", err)
		return
	}
}

func TestService_Pay_Success(t *testing.T) {
	s := newTestService()
	phone := types.Phone("+992927894561")
	amountAccount := types.Money(10)
	amountPay := types.Money(9)
	categoryPay := types.PaymentCategory("auto")
	account, err := s.addAccountWithBalance(phone, amountAccount)

	if err != nil {
		t.Errorf("Pay(): can't pay account, error %v", err)
		return
	}

	gotPayment, err := s.Pay(account.ID, amountPay, categoryPay)

	if err != nil {
		t.Errorf("Pay(): can't Pay account, error %v", err)
		return
	}

	payment := &types.Payment{
		AccountID: account.ID,
		Amount:    amountPay,
		Category:  categoryPay,
		Status:    types.PaymentStatusInProgress,
	}
	if !reflect.DeepEqual(payment.GetPayment(), gotPayment.GetPayment()) {
		t.Errorf("Pay(): can't Pay account, error %v", err)
		return
	}

	if !reflect.DeepEqual(account.Balance, amountAccount-amountPay) {
		t.Errorf("Pay(): can't Pay account, error %v", err)
		return
	}
}

func TestService_Pay_AmountNegative(t *testing.T) {
	s := newTestService()
	phone := types.Phone("+992927894561")
	amountAccount := types.Money(10)
	amountPay := types.Money(-9)
	categoryPay := types.PaymentCategory("auto")
	account, err := s.addAccountWithBalance(phone, amountAccount)

	if err != nil {
		t.Errorf("Pay(): can't pay account, error %v", err)
		return
	}

	_, err = s.Pay(account.ID, amountPay, categoryPay)

	if err == nil {
		t.Errorf("Pay(): can't pay account, error %v", err)
		return
	}
}

func TestService_Pay_CannotAccount(t *testing.T) {
	s := newTestService()
	amountPay := types.Money(9)
	categoryPay := types.PaymentCategory("auto")

	_, err := s.Pay(2, amountPay, categoryPay)

	if err == nil {
		t.Errorf("Pay(): can't pay account, error %v", err)
		return
	}
}

func TestService_Pay_NotEnoughBalance(t *testing.T) {
	s := newTestService()
	phone := types.Phone("+992927894561")
	amountAccount := types.Money(10)
	amountPay := types.Money(11)
	categoryPay := types.PaymentCategory("auto")
	account, err := s.addAccountWithBalance(phone, amountAccount)

	if err != nil {
		t.Errorf("Pay(): can't pay account, error %v", err)
		return
	}

	_, err = s.Pay(account.ID, amountPay, categoryPay)

	if err == nil {
		t.Errorf("Pay(): can't pay account, error %v", err)
		return
	}
}

func TestService_FindAccountByID_success(t *testing.T) {
	s := newTestService()
	account, _, err := s.addDefaultAccountAndPayment()

	if err != nil {
		t.Errorf("FindAccountByID(): can't find account, error %v", err)
		return
	}

	gotAccount, err := s.FindAccountByID(account.ID)

	if err != nil {
		t.Errorf("FindAccountByID(): can't find account, error %v", err)
		return
	}

	if !reflect.DeepEqual(account, gotAccount) {
		t.Errorf("FindAccountByID(): can't find account, error %v", err)
		return
	}
}

func TestService_FindAccountByID_error(t *testing.T) {
	s := newTestService()
	_, _, err := s.addDefaultAccountAndPayment()

	if err != nil {
		t.Errorf("FindAccountByID(): can't find account, error %v", err)
		return
	}

	_, err = s.FindAccountByID(2)
	if err == nil {
		t.Errorf("FindAccountByID(): can't find account, error %v", err)
		return
	}

}

func TestService_FindPaymentByID_success(t *testing.T) {
	s := newTestService()
	_, gotPayment, err := s.addDefaultAccountAndPayment()

	if err != nil {
		t.Errorf("FindPaymentByID(): can't find payment, error %v", err)
		return
	}

	resultPayment, err := s.FindPaymentByID(gotPayment.ID)

	if err != nil {
		t.Errorf("FindPaymentByID(): can't find payment, error %v", err)
		return
	}

	if !reflect.DeepEqual(resultPayment.GetPayment(), gotPayment.GetPayment()) {
		t.Errorf("Pay(): can't Pay account, error %v", err)
		return
	}
}

func TestService_FindPaymentByID_error(t *testing.T) {
	s := newTestService()
	_, _, err := s.addDefaultAccountAndPayment()

	if err != nil {
		t.Errorf("FindPaymentByID(): can't find payment, error %v", err)
		return
	}

	_, err = s.FindPaymentByID("tes")

	if err == nil {
		t.Errorf("FindPaymentByID(): can't find payment, error %v", err)
		return
	}
}

func TestService_Reject_success(t *testing.T) {
	s := newTestService()
	_, gotPayment, err := s.addDefaultAccountAndPayment()

	if err != nil {
		t.Errorf("Reject(): can't Reject payment, error %v", err)
		return
	}

	err = s.Reject(gotPayment.ID)

	if err != nil || gotPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): can't Reject payment, error %v", err)
		return
	}

}

func TestService_Reject_NotFountPayment(t *testing.T) {
	s := newTestService()
	_, _, err := s.addDefaultAccountAndPayment()

	if err != nil {
		t.Errorf("Reject(): can't Reject payment, error %v", err)
		return
	}

	err = s.Reject("")

	if err == nil {
		t.Errorf("Reject(): can't Reject payment, error %v", err)
		return
	}

}

func TestService_Reject_NotFountAccount(t *testing.T) {
	s := newTestService()
	_, gotPayment, err := s.addDefaultAccountAndPayment()

	if err != nil {
		t.Errorf("Reject(): can't Reject payment, error %v", err)
		return
	}

	gotPayment.AccountID = 3
	err = s.Reject(gotPayment.ID)

	if err == nil {
		t.Errorf("Reject(): can't Reject payment, error %v", err)
		return
	}

}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	_, gotPayment, err := s.addDefaultAccountAndPayment()

	if err != nil {
		t.Errorf("Repeat(): can't Repeat payment, error %v", err)
		return
	}

	result, err := s.Repeat(gotPayment.ID)

	if !reflect.DeepEqual(result.GetPayment(), gotPayment.GetPayment()) {
		t.Errorf("Repeat(): can't Repeat account, error %v", err)
		return
	}
}

func TestService_Repeat_error(t *testing.T) {
	s := newTestService()
	_, err := s.Repeat("gotPayment.ID")

	if err == nil {
		t.Errorf("Repeat(): can't Repeat account, error %v", err)
		return
	}
}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()
	_, gotPayment, err := s.addDefaultAccountAndPayment()

	if err != nil {
		t.Errorf("FavoritePayment(): can't Favorite payment, error %v", err)
		return
	}

	_, err = s.FavoritePayment(gotPayment.ID, "test")

	if err != nil {
		t.Errorf("FavoritePayment(): can't Favorite Payment, error %v", err)
		return
	}
}

func TestService_FavoritePayment_error(t *testing.T) {
	s := newTestService()
	_, err := s.FavoritePayment("gotPayment.ID", "test")

	if err == nil {
		t.Errorf("Repeat(): can't Repeat account, error %v", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()
	_, gotPayment, err := s.addDefaultAccountAndPayment()

	if err != nil {
		t.Errorf("PayFromFavorite(): can't PayFromFavorite payment, error %v", err)
		return
	}

	favorite, err := s.FavoritePayment(gotPayment.ID, "test")

	if err != nil {
		t.Errorf("PayFromFavorite(): can't PayFromFavorite Payment, error %v", err)
		return
	}

	payment, err := s.PayFromFavorite(favorite.ID)

	if err != nil {
		t.Errorf("PayFromFavorite(): can't PayFromFavorite Payment, error %v", err)
		return
	}

	result, err := s.FindPaymentByID(payment.ID)

	if err != nil || !reflect.DeepEqual(result.GetPayment(), payment.GetPayment()) {
		t.Errorf("PayFromFavorite(): can't PayFromFavorite Payment, error %v", err)
		return
	}
}

func TestService_PayFromFavorite_error(t *testing.T) {
	s := newTestService()
	_, err := s.PayFromFavorite("gotPayment.ID")

	if err == nil {
		t.Errorf("PayFromFavorite(): can't PayFromFavorite payment, error %v", err)
		return
	}
}
