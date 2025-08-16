package model

import "github.com/asaskevich/govalidator"

type NewUser struct {
	Name        string  `json:"name" valid:"required"`
	Email       string  `json:"email" valid:"required"`
	PhoneNumber string  `json:"phoneNumber" valid:"required"`
	Line1       string  `json:"line1" valid:"required"`
	Line2       *string `json:"line2"`
	Line3       *string `json:"line3"`
	Town        string  `json:"town" valid:"required"`
	County      *string `json:"county"`
	Postcode    string  `json:"postcode" valid:"required"`
}

type User struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	Password    string  `json:"-"` // never expose password
	Status      string  `json:"status"`
	Line1       string  `json:"line1"`
	Line2       *string `json:"line2"`
	Line3       *string `json:"line3"`
	Town        string  `json:"town"`
	County      *string `json:"county"`
	Postcode    string  `json:"postcode"`
}

func (n *NewUser) Valid() (bool, error) {
	return govalidator.ValidateStruct(n)
}
