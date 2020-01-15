package payment
import (
	"time"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type OmisePayment struct {
	client *omise.Client
}

func NewOmisePayment(publicKey string, secretKey string) (*OmisePayment, error) {
	c, e := omise.NewClient(publicKey, secretKey)
	return &OmisePayment{client: c}, e
}

func (o *OmisePayment) CreateToken(name string, number string, expMonth time.Month, expYear int, cvv string) (*omise.Token, error){
	token, createToken := &omise.Token{}, &operations.CreateToken{
		Name:           name,
		Number:         number,
		ExpirationMonth: expMonth,
		ExpirationYear:  expYear,
		SecurityCode:  cvv,
	}
  	e := o.client.Do(token, createToken)
  	return token, e
}

func (o *OmisePayment) Charge(amount int64, currency string, token *omise.Token) (*omise.Charge, error){
	charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
		Amount:   amount, 
		Currency: currency,
		Card:     token.ID,
	}
	e := o.client.Do(charge, createCharge)
	return charge, e
}