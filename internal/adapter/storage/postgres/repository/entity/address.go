package entity

import (
	"time"

	"eagle-bank.com/internal/core/domain/model"
	"github.com/google/uuid"
)

func NewAddress(opts ...Option[*Address]) (Address, error) {
	newEntity := Address{
		id:        ID(uuid.NewString()),
		createdAt: time.Now().UTC(),
	}
	err := newEntity.Modify(opts...)
	if err != nil {
		return Address{}, err
	}
	return newEntity, nil
}

func (a *Address) Modify(opts ...Option[*Address]) error {
	cl, err := Clone(a)
	if err != nil {
		return err
	}
	ApplyOptions(opts, cl)

	err = validate(userAddressValidation{
		ID:        cl.id,
		UserID:    cl.userID,
		Line1:     cl.line1,
		Town:      cl.town,
		Postcode:  cl.postcode,
		CreatedAt: cl.createdAt,
	})
	if err != nil {
		return err
	}

	*a = *cl
	return nil
}

type Address struct {
	id        ID
	userID    ID
	line1     string
	line2     *string
	line3     *string
	town      string
	county    *string
	postcode  string
	createdAt time.Time
	updatedAt time.Time
}

type userAddressValidation struct {
	ID        ID        `valid:"uuid,required"`
	UserID    ID        `valid:"uuid,required"`
	Line1     string    `valid:"required"`
	Town      string    `valid:"required"`
	Postcode  string    `valid:"required"`
	CreatedAt time.Time `valid:"required"`
}

func (a *Address) ID() ID {
	return a.id
}

func (a *Address) UserID() ID {
	return a.userID
}

func (a *Address) Line1() string {
	return a.line1
}

func (a *Address) Line2() *string {
	return a.line2
}

func (a *Address) Line3() *string {
	return a.line2
}

func (a *Address) Town() string {
	return a.town
}

func (a *Address) County() *string {
	return a.county
}

func (a *Address) Postcode() string {
	return a.postcode
}

func (a *Address) CreatedAt() time.Time {
	return a.createdAt
}

func WithUserAddressID(id ID) Option[*Address] {
	return func(a *Address) {
		a.id = id
	}
}

func WithUserAddressUserID(id ID) Option[*Address] {
	return func(a *Address) {
		a.userID = id
	}
}

func WithUserAddressLine1(line1 string) Option[*Address] {
	return func(a *Address) {
		a.line1 = line1
	}
}

func WithUserAddressLine2(line2 *string) Option[*Address] {
	return func(a *Address) {
		a.line2 = line2
	}
}

func WithUserAddressLine3(line3 *string) Option[*Address] {
	return func(a *Address) {
		a.line3 = line3
	}
}

func WithUserAddressTown(town string) Option[*Address] {
	return func(a *Address) {
		a.town = town
	}
}

func WithUserAddressCounty(county *string) Option[*Address] {
	return func(a *Address) {
		a.county = county
	}
}

func WithUserAddressPostcode(postcode string) Option[*Address] {
	return func(a *Address) {
		a.postcode = postcode
	}
}

func WithUserAddressCreatedAt(createdAt time.Time) Option[*Address] {
	return func(a *Address) {
		a.createdAt = createdAt
	}
}

type AddressDAO struct {
	ID        ID        `db:"id"`
	UserID    ID        `db:"user_id"`
	Line1     string    `db:"line1"`
	Line2     *string   `db:"line2"`
	Line3     *string   `db:"line3"`
	Town      string    `db:"town"`
	County    *string   `db:"county"`
	Postcode  string    `db:"postcode"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (a *Address) FromEntity() AddressDAO {
	return AddressDAO{
		ID:        a.id,
		UserID:    a.userID,
		Line1:     a.line1,
		Line2:     a.line2,
		Line3:     a.line3,
		Town:      a.town,
		County:    a.county,
		Postcode:  a.postcode,
		CreatedAt: a.createdAt,
		UpdatedAt: a.updatedAt,
	}
}

func ConvertUserAddressFromModel(m *model.User) *AddressDAO {
	return &AddressDAO{
		Line1:    m.Line1,
		Line2:    m.Line2,
		Line3:    m.Line3,
		Town:     m.Town,
		County:   m.County,
		Postcode: m.Postcode,
	}
}
