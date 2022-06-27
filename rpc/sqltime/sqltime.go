// Extended from: https://github.com/SamuelTissot/sqltime

package sqltime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// the degree of precision to REMOVE
// Default time.Microsecond
var TruncateOff = time.Microsecond

// Local the timezone the database is set to
// default UTC
var DatabaseLocation, _ = time.LoadLocation("UTC")

// Time
// type that can be used with sql driver's and offers
// a less precise sql timestamp
type Time struct {
	time.Time
}

// satisfy the sql.scanner interface
func (t *Time) Scan(value interface{}) error {
	rt, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("dbtime could not convert value into time.Time. value: %v", value)
	}
	*t = Time{format(rt)}
	return nil
}

// satifies the driver.Value interface
func (t Time) Value() (driver.Value, error) {
	return format(t.Time), nil //format just in case
}

// Now wrapper around the time.Now() function
func Now() Time {
	return Time{format(time.Now())}
}

// Date wrapper around the time.Date() function
func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Time {
	return Time{format(time.Date(year, month, day, hour, min, sec, nsec, loc))}
}

// insure the correct format
func format(t time.Time) time.Time {
	return t.In(DatabaseLocation).Truncate(TruncateOff)
}

/*
	Extended function support to convert classic golang time.Time to SQL time.
*/

func New(t time.Time) Time {
	return Time{format(t)}
}
