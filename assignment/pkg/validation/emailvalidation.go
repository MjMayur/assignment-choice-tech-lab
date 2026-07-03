package validation

import (
	"regexp"
	"strings"
)

func EmailValidation(email string) bool {

	email = strings.TrimSpace(email)

	expression := `^[a-zA-Z0-9]+([._%+-]?[a-zA-Z0-9]+)*@([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]{2,}$`

	// Stricter regex:
	// 1. Local part: letters, numbers, ._%+-
	// 2. Domain part: must not start or end with '-'
	// 3. TLD: at least 2 letters
	// expression := `^[a-zA-Z0-9]+([._%+-]?[a-zA-Z0-9]+)*@([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]{2,}$`

	re := regexp.MustCompile(expression)
	return re.MatchString(email)
}
