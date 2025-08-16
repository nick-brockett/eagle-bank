package entity

import (
	"github.com/asaskevich/govalidator"
)

func validate(value any) error {
	valid, err := govalidator.ValidateStruct(value)

	if !valid {
		return err
	}
	return nil
}
