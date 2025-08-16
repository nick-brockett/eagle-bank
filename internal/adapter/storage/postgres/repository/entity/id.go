package entity

type ID string

func (i ID) String() string {
	return string(i)
}
