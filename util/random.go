package util

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// emptyLocation creates a timezone location that matches PostgreSQL's behavior
// This matches exactly what the PostgreSQL driver returns
var emptyLocation = time.FixedZone("", 0)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomName generates a random name
func RandomName() string {
	return RandomString(6)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@%s.com", RandomString(6), RandomString(4))
}

// RandomPhone generates a random phone number
func RandomPhone() string {
	return fmt.Sprintf("+1%d", RandomInt(1000000000, 9999999999))
}

// RandomPrice generates a random price string
func RandomPrice() string {
	return fmt.Sprintf("%.2f", float64(RandomInt(100, 10000))/100)
}

// RandomWeek generates a random week number
func RandomWeek() int32 {
	return int32(RandomInt(1, 52))
}

// RandomDate generates a random date within the last year
func RandomDate() time.Time {
	now := time.Now()
	oneYearAgo := now.AddDate(-1, 0, 0)
	delta := now.Sub(oneYearAgo)
	randomDuration := time.Duration(rand.Int63n(int64(delta)))
	randomDate := oneYearAgo.Add(randomDuration)
	// Return only the date part with the same timezone as database
	return time.Date(randomDate.Year(), randomDate.Month(), randomDate.Day(), 0, 0, 0, 0, emptyLocation)
}

// RandomTime generates a random time of day
func RandomTime() time.Time {
	hour := rand.Intn(24)
	minute := rand.Intn(60)
	second := rand.Intn(60)
	// Use epoch date (1970-01-01) for time-only fields with same timezone as database
	return time.Date(1970, 1, 1, hour, minute, second, 0, emptyLocation)
}

// NullString creates a sql.NullString from a string
func NullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

// NullInt32 creates a sql.NullInt32 from an int32
func NullInt32(i int32) sql.NullInt32 {
	return sql.NullInt32{Int32: i, Valid: true}
}

// NullTime creates a sql.NullTime from a time.Time
func NullTime(t time.Time) sql.NullTime {
	return sql.NullTime{Time: t, Valid: true}
}

// RandomBool generates a random boolean
func RandomBool() bool {
	return rand.Intn(2) == 1
}

// RandomFloat generates a random float between min and max
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// RandomMoney generates a random money amount as string
func RandomMoney() string {
	return fmt.Sprintf("%.2f", RandomFloat(0.01, 999.99))
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "GBP", "CAD", "AUD", "JPY", "INR"}
	return currencies[rand.Intn(len(currencies))]
}

// RandomPhoneWithFormat generates a phone number with specific format
func RandomPhoneWithFormat() string {
	return fmt.Sprintf("(%03d) %03d-%04d",
		RandomInt(100, 999),
		RandomInt(100, 999),
		RandomInt(1000, 9999))
}

// TimeEqual compares two sql.NullTime values, ignoring timezone and date for time-only comparisons
func TimeEqual(expected, actual sql.NullTime) bool {
	if expected.Valid != actual.Valid {
		return false
	}
	if !expected.Valid {
		return true
	}

	e, a := expected.Time, actual.Time

	// If either has year 0 or 1970, compare only time components (for PostgreSQL time type)
	if (e.Year() == 0 || e.Year() == 1970) && (a.Year() == 0 || a.Year() == 1970) {
		return e.Hour() == a.Hour() && e.Minute() == a.Minute() && e.Second() == a.Second()
	}

	// For date fields, compare year, month, day
	return e.Year() == a.Year() && e.Month() == a.Month() && e.Day() == a.Day()
}
