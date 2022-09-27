package pixie

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	cru2 "github.com/tuupke/pixie/crud"
)

type Settings gorm.DB

func (set *Settings) CrudController() *cru2.Controller[Setting] {
	return cru2.New[Setting]((*gorm.DB)(set))
}

type RawMessage []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil || string(m) == "null" {
		return []byte("null"), nil
	}

	if len(m) == 0 {
		return m, nil
	}

	// check if rawMessage is a json, if not pre and append a '"'
	fc := string(m[0])
	if fc == "[" || fc == "{" {
		return m, nil
	}

	return append(append([]byte("\""), m...), byte('"')), nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

func (m *RawMessage) Scan(src any) error {
	if src == nil {
		*m = nil
		return nil
	}

	switch v := src.(type) {
	case []byte:
		*m = v
	case string:
		*m = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal raw: %w", src)
	}

	return nil
}

type Setting struct {
	Key   string     `gorm:"primaryKey" json:"key"`
	Value RawMessage `json:"value"`
}

func (s Setting) String() string {
	return s.Key
}

func (s Setting) Identifier() fmt.Stringer {
	return s
}

// var settings = new(Settings)

func LoadSettings(orm *gorm.DB) *Settings {
	return (*Settings)(orm)
}

func (s *Settings) Retrieve(k string) (v string) {
	v, _ = s.Get(k)
	return
}

func (s *Settings) Fallback(k, fb string) (v string) {
	var found bool
	if v, found = s.Get(k); !found {
		v = fb
	}
	return
}

func (s *Settings) Has(k string) (exists bool) {
	var setting []Setting
	((*gorm.DB)(s)).Model(Setting{}).Where(clause.IN{
		Column: clause.Column{
			Table: clause.CurrentTable,
			Name:  "key",
			Raw:   false,
		},
		Values: []interface{}{k},
	}).Find(&setting)

	return len(setting) > 0
}

func (s *Settings) Get(k string) (v string, found bool) {
	var setting []Setting
	((*gorm.DB)(s)).Model(Setting{}).Where(clause.IN{
		Column: clause.Column{
			Table: clause.CurrentTable,
			Name:  "key",
			Raw:   false,
		},
		Values: []interface{}{k},
	}).Find(&setting)

	if len(setting) > 0 {
		v = string(setting[0].Value)
	}

	return
}

func (s *Settings) Set(k string, v []byte) {
	setting := &Setting{
		Key:   k,
		Value: v,
	}

	(*gorm.DB)(s).Model(setting).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(setting)
}

func (s Settings) SetMultiple(settings ...Setting) {
	for _, setting := range settings {
		s.Set(setting.Key, setting.Value)
	}
}
