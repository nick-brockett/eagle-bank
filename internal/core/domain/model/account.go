package model

import "time"

type NewAccount struct {
	UserID        string `json:"userId" valid:"required"`
	Name          string `json:"name" valid:"required"`
	Type          string `json:"type" valid:"required"`
	AccountNumber string `json:"-"`
}

type Account struct {
	AccountNumber    string    `json:"accountNumber"`
	SortCode         string    `json:"sortCode"`
	Name             string    `json:"name"`
	AccountType      string    `json:"accountType"`
	Balance          int       `json:"balance"`
	Currency         string    `json:"currency"`
	CreatedTimestamp time.Time `json:"createdTimestamp"`
	UpdatedTimestamp time.Time `json:"updatedTimestamp"`
}

type UserAccount struct {
	UserID        string `json:"userId"`
	AccountNumber string `json:"accountNumber"`
}
