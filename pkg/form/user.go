package form

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailReg ...
var EmailReg = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// UserForm ...
type UserForm struct {
	url.Values        // hold form data
	Errors     errors // validation errors
}

// Init initialize a user form struct
func (f *UserForm) Init(v url.Values) *UserForm {
	return &UserForm{
		v,
		errors(map[string][]string{}),
	}
}

// Required check fields are not blank
func (f *UserForm) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "Required")
		}
	}
}

// MaxLength check field string maximum length
func (f *UserForm) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("Max string length is %d", d))
	}
}

// MinLength check field string minimum length
func (f *UserForm) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("Min string length is %d", d))
	}
}

// Match match string with regexp
func (f *UserForm) Match(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "Invalid")
	}
}

// Valid method returns true if there are no errors.
func (f *UserForm) Valid() bool {
	return len(f.Errors) == 0
}
