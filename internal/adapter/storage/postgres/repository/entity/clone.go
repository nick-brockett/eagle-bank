package entity

import (
	"github.com/jinzhu/copier"
)

// Clone makes a deep copy of a pointer to a struct
func Clone[T any](in *T) (*T, error) {
	out := new(T)
	err := copier.Copy(out, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}
