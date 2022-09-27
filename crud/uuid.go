package crud

import (
	"database/sql/driver"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type (
	UUID uuid.UUID
)

func (u UUID) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	oType, oSize := field.DataType, field.Size
	defer func() { field.DataType = oType; field.Size = oSize }()

	field.DataType = schema.String

	field.Size = 36

	return db.Dialector.DataTypeOf(field)
}

func NewUUID() UUID {
	return UUID(uuid.New())
}

func (u *UUID) Scan(src interface{}) error {
	return (*uuid.UUID)(u).Scan(src)
}

func (u UUID) Value() (driver.Value, error) {
	if u.Blank() {
		return nil, nil
	}
	return (uuid.UUID)(u).Value()
}

func (u UUID) String() string {
	return (uuid.UUID)(u).String()
}

func UUIDFromString(s string) (UUID, error) {
	u, err := uuid.Parse(s)
	return UUID(u), err
}

func (u UUID) Blank() bool {
	return u == UUID(uuid.Nil)
}

// MarshalText implements encoding.TextMarshaler.
func (u UUID) MarshalText() ([]byte, error) {
	if u == UUID(uuid.Nil) {
		return nil, nil
	}

	return (uuid.UUID)(u).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *UUID) UnmarshalText(data []byte) error {
	return (*uuid.UUID)(u).UnmarshalText(data)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (u UUID) MarshalBinary() ([]byte, error) {
	return (uuid.UUID)(u).MarshalBinary()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (u *UUID) UnmarshalBinary(data []byte) error {
	return (*uuid.UUID)(u).UnmarshalBinary(data)
}
