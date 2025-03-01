package utils

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

// TemplateFuncs returns a map of template functions
func TemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatTime": formatTime,
		"dict":       dict,
		"truncate":   truncate,
		"join":       strings.Join,
		"contains":   contains,
		"len":        length,
	}
}

// formatTime formats a time.Time to a readable string
func formatTime(t time.Time) string {
	return t.Format("Jan 02, 2006 15:04")
}

// dict creates a map from a list of key/value pairs
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid dict call, must have even number of args")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// truncate cuts a string to the specified length and adds ellipsis if needed
func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// contains checks if a string is present in a slice
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// length is a template function for len()
func length(v interface{}) int {
	switch val := v.(type) {
	case []interface{}:
		return len(val)
	case map[string]interface{}:
		return len(val)
	case string:
		return len(val)
	default:
		return 0
	}
}