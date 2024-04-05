package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Form struct {
	Errors map[string][]string
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func New() *Form {
	return &Form{make(map[string][]string)}
}

func (f *Form) EmailValid(email string) bool {
	return emailRegex.Match([]byte(email))
}

func (f *Form) NameValid(name string) bool {
	return nameRegex.Match([]byte(name))
}

// Должен содержать как минимум 8 символов, включая по крайней мере одну заглавную букву,
// одну строчную букву, одну цифру и один специальный символ.

func (f *Form) PasswordValid(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasUpperCase bool
		hasLowerCase bool
		hasDigit     bool
		hasSpecial   bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpperCase = true
		case unicode.IsLower(char):
			hasLowerCase = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpperCase && hasLowerCase && hasDigit && hasSpecial
}

func (f *Form) EmptyFieldValid(val string) bool {
	return strings.TrimSpace(val) != ""
}

func (f *Form) MinLengthValid(val string, length int) bool {
	return len(val) >= length
}

var (
	nameRegex  = regexp.MustCompile("^[a-zA-Zа-яА-Я0-9_]{3,16}$")
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)


func (f *Form) CheckField(ok bool, field, err string) error {
	if !ok {
		f.Errors[field] = append(f.Errors[field], err)
		return fmt.Errorf(err)
	}
	return nil
}
