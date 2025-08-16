package dao

import (
	"eagle-bank.com/internal/core/domain/model"
)

type UserViewDAO struct {
	ID          string  `db:"id"`
	Name        string  `db:"name"`
	Email       string  `db:"email"`
	PhoneNumber string  `db:"phone_number"`
	Status      string  `db:"status"`
	Line1       string  `db:"line1"`
	Line2       *string `db:"line2"`
	Line3       *string `db:"line3"`
	Town        string  `db:"town"`
	County      *string `db:"county"`
	Postcode    string  `db:"postcode"`
}

func (u UserViewDAO) ConvertToModel() *model.User {
	return &model.User{
		ID:          u.ID,
		Name:        u.Name,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
		Status:      u.Status,
		Line1:       u.Line1,
		Line2:       u.Line2,
		Line3:       u.Line3,
		Town:        u.Town,
		County:      u.County,
		Postcode:    u.Postcode,
	}
}
