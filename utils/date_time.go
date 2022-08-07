package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	BatchFmt    = "20060102"
	DateFmt     = "2006-01-02"
	DateTimeFmt = "2006-01-02 15:04:05"
)

type DateTime struct {
	time.Time
}

type Date struct {
	time.Time
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("\"%v\"", t.Format(DateTimeFmt))
	return []byte(s), nil
}

func (t *DateTime) UnmarshalJSON(data []byte) (err error) {
	t.Time, _ = time.ParseInLocation(`"`+DateTimeFmt+`"`, string(data), time.Local)
	return
}

func (t *DateTime) String() string {
	str, _ := json.Marshal(t)
	return string(str)
}

func (t *DateTime) SetRaw(value interface{}) error {
	switch value.(type) {
	case time.Time:
		t.Time = value.(time.Time)
	}
	return nil
}

func (t *DateTime) RawValue() interface{} {
	str := t.Format(DateTimeFmt)
	if str == "0001-01-01 00:00：00" {
		return nil
	}
	return str
}

func (t Date) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf("\"%v\"", t.Format(DateFmt))
	return []byte(s), nil
}

func (t *Date) UnmarshalJSON(data []byte) (err error) {
	t.Time, _ = time.ParseInLocation(`"`+DateFmt+`"`, string(data), time.Local)
	return
}

func (t *Date) String() string {
	str, _ := json.Marshal(t)
	return string(str)
}

func (t *Date) SetRaw(value interface{}) error {
	switch value.(type) {
	case time.Time:
		t.Time = value.(time.Time)
	}
	return nil
}

func (t *Date) RawValue() interface{} {
	str := t.Format(DateFmt)
	if str == "0001-01-01" {
		return nil
	}
	return str
}

func TimeDateAfter(t0 time.Time, t1 time.Time) bool {
	if t0.Year() < t1.Year() {
		return false
	} else if t0.Year() > t1.Year() {
		return true
	}
	return t0.YearDay() > t1.YearDay()
}

func TimeDateEqual(t0 time.Time, t1 time.Time) bool {
	return ((t0.Year() == t1.Year()) && t0.YearDay() == t1.YearDay())
}
