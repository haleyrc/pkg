package time

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const EpochTime = "2022-09-16T15:02:04Z"

type timeFn func() time.Time

var nowFn timeFn = time.Now

var (
	Hour = Duration{d: time.Hour}
)

type Duration struct {
	d time.Duration
}

type Time struct {
	t time.Time
}

func (t Time) Add(d Duration) Time {
	return Time{t: t.t.Add(d.d).Truncate(time.Second)}
}

func (t Time) Before(other Time) bool {
	return t.t.Before(other.t)
}

func (t Time) Equal(other Time) bool {
	return t.t.Equal(other.t)
}

func (t Time) IsZero() bool {
	return t.t.IsZero()
}

func (t Time) MarshalJSON() ([]byte, error) {
	return t.t.MarshalJSON()
}

func (t *Time) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch value := value.(type) {
	case time.Time:
		t.t = value.Truncate(time.Second)
	default:
		return fmt.Errorf("unsupported data type: %T", value)
	}
	return nil
}

func (t Time) String() string {
	return t.t.String()
}

func (t Time) Sub(d Duration) Time {
	return Time{t: t.t.Add(-d.d).Truncate(time.Second)}
}

func (t Time) UTC() Time {
	return Time{t: t.t.UTC().Truncate(time.Second)}
}

func (t Time) Value() (driver.Value, error) {
	return t.t, nil
}

func Freeze() {
	nowFn = func() time.Time {
		t, err := time.ParseInLocation(time.RFC3339, EpochTime, time.UTC)
		if err != nil {
			panic(err)
		}
		return t
	}
}

func Unfreeze() {
	nowFn = time.Now
}

func Now() Time {
	return Time{t: nowFn().Truncate(time.Second)}
}
