package templates

import (
	"strings"
	"text/template"
)

// GetTemplateFuncs returns template functions
func GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"lower":          LowerFirstChar,
		"upper":          strings.ToUpper,
		"lowerFirstChar": LowerFirstChar,
		"snakeCase":      ToSnakeCase,
		"contains":       strings.Contains, // Added contains function for template conditionals
	}
}

func ToSnakeCase(s string) string {
	result := ""
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}

func LowerFirstChar(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
