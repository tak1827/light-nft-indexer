package util

import (
	"errors"
	"fmt"
	"math"
	"net/mail"
	"reflect"
	"regexp"
	"unicode"
)

var (
	regNumber      = regexp.MustCompile(`^[0-9]+$`)
	regAlpha       = regexp.MustCompile(`^[a-zA-Z]+$`)
	regAlphaNumber = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	regAscii       = regexp.MustCompile(`[a-zA-Z0-9!-/:-@¥[-{-~]+$`)
	regPhoneNumber = regexp.MustCompile(`^0\d{2,3}-\d{1,4}-\d{4}$`)

	ErrNotNumber         = errors.New("not number")
	ErrInvalidEthAddress = errors.New("invalid eth address")
)

func ValidateOnlyNumber(src string) error {
	if !regNumber.MatchString(src) {
		return ErrNotNumber
	}
	return nil
}

// TODO: should write test code
func ValidateOnlyAscii(src string) error {
	if !regAscii.MatchString(src) {
		return errors.New("only ascii are allowed")
	}
	return nil
}

func ValidateOnlyAlphaNumber(src string) error {
	if !regAlphaNumber.MatchString(src) {
		return errors.New("only alphabet and number are allowed")
	}
	return nil
}

// NOTE: is greater than or equal to `>=`
func ValidateGE(min int, src interface{}) error {
	var (
		v      = reflect.ValueOf(src)
		length int
	)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice:
		length = v.Len()
	case reflect.Uint64:
		length = int(v.Uint())
	case reflect.Float64:
		length = int(math.Floor(v.Float())) // -0.1 -> -1
	default:
		return errors.New("unexpected reflect.Kind")
	}

	if length < min {
		return errors.New(fmt.Sprintf("length(=%d) should be equal or greater than %d", length, min))
	}

	return nil
}

// NOTE: is less than or equal to `<=`
func ValidateLE(max int, src interface{}) error {
	var (
		v      = reflect.ValueOf(src)
		length int
	)
	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice:
		length = v.Len()
	case reflect.Uint64:
		length = int(v.Uint())
	case reflect.Float64:
		length = int(math.Ceil(v.Float())) // 100.1 -> 101
	default:
		return errors.New("unexpected reflect.Kind")
	}

	if length > max {
		return errors.New(fmt.Sprintf("length(=%d) should be equal or less than %d", length, max))
	}

	return nil
}

func ValidateLenBetween(min, max int, src interface{}) error {
	if err := ValidateGE(min, src); err != nil {
		return err
	}
	return ValidateLE(max, src)
}

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("at mail.ParseAddress: %w", err)
	}
	return nil
}

func ValidateTelNum(telNum string) error {
	if !regPhoneNumber.MatchString(telNum) {
		return errors.New("invalid telNum")
	}
	return nil
}

// NOTE: Ref
// http://www.inanzzz.com/index.php/post/8l1a/validating-user-password-in-golang-requests
func ValidatePassword(pass string, min int) error {
	var (
		upp, low, num, sym bool
	)

	if len(pass) < min {
		return fmt.Errorf("password should be longer or equal than %d", min)
	}

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upp = true
		case unicode.IsLower(char):
			low = true
		case unicode.IsNumber(char):
			num = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
		default:
			return errors.New("invalid letter is in password")
		}
	}

	if !upp {
		return errors.New("should include upper case lette at lease once in password")
	}

	if !low {
		return errors.New("should include lower case lette at lease once in password")
	}

	if !num {
		return errors.New("should include lower case lette at lease once in password")
	}

	if !sym {
		return errors.New("should include symbol at lease once in password")
	}

	return nil
}

func ValidateEthAddress(address string) error {
	if address[:2] != "0x" || len(address[2:]) != 40 {
		return ErrInvalidEthAddress
	}

	return nil
}
