package service

import "time"

// CalculateAge computes the exact age based on the current date.
func CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()

	// If the current month is before the birth month, or if it's the birth month 
	// but the current day is before the birth day, they haven't had their birthday yet this year.
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		age--
	}

	return age
}