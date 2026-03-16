package validation

import (
	"net/mail"
	"net/url"
	"strings"
	"time"
)

type Errors map[string]string

func NewErrors() Errors {
	return Errors{}
}

func (e Errors) Add(field string, message string) {
	field = strings.TrimSpace(field)
	message = strings.TrimSpace(message)
	if field == "" || message == "" {
		return
	}

	if _, exists := e[field]; !exists {
		e[field] = message
	}
}

func (e Errors) HasAny() bool {
	return len(e) > 0
}

func (e Errors) RequiredString(field string, value string, message string) {
	if strings.TrimSpace(value) == "" {
		e.Add(field, message)
	}
}

func (e Errors) OptionalMaxLength(field string, value *string, max int, message string) {
	if value == nil {
		return
	}

	if len(strings.TrimSpace(*value)) > max {
		e.Add(field, message)
	}
}

func (e Errors) RequiredMaxLength(field string, value string, max int, requiredMessage string, maxMessage string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		e.Add(field, requiredMessage)
		return
	}

	if len(trimmed) > max {
		e.Add(field, maxMessage)
	}
}

func (e Errors) Email(field string, value string, message string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return
	}

	if _, err := mail.ParseAddress(trimmed); err != nil {
		e.Add(field, message)
	}
}

func (e Errors) URL(field string, value *string, message string) {
	if value == nil {
		return
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return
	}

	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		e.Add(field, message)
	}
}

func (e Errors) Date(field string, value *string, layout string, message string) {
	if value == nil {
		return
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return
	}

	if _, err := time.Parse(layout, trimmed); err != nil {
		e.Add(field, message)
	}
}

func (e Errors) Enum(field string, value string, allowed []string, message string) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return
	}

	for _, option := range allowed {
		if trimmed == option {
			return
		}
	}

	e.Add(field, message)
}

func (e Errors) MinLength(field string, value string, min int, message string) {
	if len(strings.TrimSpace(value)) < min {
		e.Add(field, message)
	}
}

func (e Errors) PositiveInt64(field string, value int64, message string) {
	if value <= 0 {
		e.Add(field, message)
	}
}

func (e Errors) PositiveInt(field string, value int, message string) {
	if value <= 0 {
		e.Add(field, message)
	}
}

func (e Errors) IntRange(field string, value int, min int, max int, message string) {
	if value < min || value > max {
		e.Add(field, message)
	}
}
