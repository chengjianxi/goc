package sqlx

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

type StringTime sql.NullTime

func StringTimeNow() StringTime {
	return StringTime{Time: time.Now(), Valid: true}
}

// Scan implements the Scanner interface.
func (n *StringTime) Scan(value interface{}) error {
	return (*sql.NullTime)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n StringTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func (n StringTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time.Format("2006-01-02 15:04:05"))
	}
	return json.Marshal("")
}

func (n *StringTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Time)
	if err == nil {
		n.Valid = true
	}
	return err
}
