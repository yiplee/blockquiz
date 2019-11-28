package property

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

const (
	timeLayout = time.RFC3339Nano
)

type Value string

func Parse(value interface{}) Value {
	switch v := value.(type) {
	case Value:
		return v
	case string:
		return Value(v)
	case []byte:
		return Value(v)
	case time.Time:
		return Value(v.Format(timeLayout))
	case int:
		return Value(strconv.Itoa(v))
	case int64:
		return Value(strconv.FormatInt(v, 10))
	default:
		return Value(fmt.Sprint(v))
	}
}

func (v Value) String() string {
	return string(v)
}

func (v Value) Time() time.Time {
	t, _ := time.Parse(timeLayout, v.String())
	return t
}

func (v Value) Int() int {
	i, _ := strconv.Atoi(v.String())
	return i
}

func (v Value) Int64() int64 {
	i, _ := strconv.ParseInt(v.String(), 10, 64)
	return i
}

func (v Value) Value() (driver.Value, error) {
	return v.String(), nil
}

func (v *Value) Scan(value interface{}) error {
	*v = Parse(value)
	return nil
}
