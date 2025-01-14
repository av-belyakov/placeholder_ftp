package supportingfunctions

import "strings"

// ReplaceCommaCharacter заменяет двойную кавычку одинарной
func ReplaceCommaCharacter(v string) string {
	return strings.ReplaceAll(v, "\"", "'")
}
