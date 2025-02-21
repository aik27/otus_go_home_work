package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrOnlyStructSupport     = errors.New("only structure type is supported")
	ErrValueNotSupported     = errors.New("value not supported")
	ErrInvalidValidationRule = errors.New("invalid validation rule")
	ErrValidateMin           = errors.New("less than min")
	ErrValidateMax           = errors.New("bigger than max")
	ErrValidateLen           = errors.New("not equal to len")
	ErrValidateIn            = errors.New("not in range")
	ErrValidateRegexp        = errors.New("does not match pattern")
)

type validationRules map[string]string

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	msg := make([]string, 0)
	for _, e := range v {
		msg = append(msg, fmt.Sprintf("%s: %v", e.Field, e.Err))
	}
	return strings.Join(msg, "\n")
}

func Validate(v interface{}) error {
	// Place your code here.
	validationErrors := make(ValidationErrors, 0)
	refValue := reflect.ValueOf(v)

	if refValue.Kind() != reflect.Struct {
		return ErrOnlyStructSupport
	}

	refType := refValue.Type()

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		value := refValue.Field(i)

		// поле публичное?
		if !refType.Field(i).IsExported() {
			continue
		}

		validateTag, ok := field.Tag.Lookup("validate")
		if !ok {
			continue
		}
		validationRules := extractValidationRules(validateTag)

		switch {
		case field.Type.Kind() == reflect.String:
			err := validateString(value.String(), validationRules)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		case field.Type.Kind() == reflect.Int:
			err := validateInt(value.Int(), validationRules)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		case field.Type.Kind() == reflect.Slice:
			err := validateSlice(value, validationRules)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: field.Name,
					Err:   err,
				})
			}
		default:
			continue
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateString(value string, rules validationRules) error {
	for key, rule := range rules {
		switch key {
		case "in":
			collection := strings.Split(rule, ",")
			if !slices.Contains(collection, value) {
				return ErrValidateIn
			}

		case "len":
			ln, err := strconv.Atoi(rule)
			if err != nil {
				return ErrInvalidValidationRule
			}

			if len([]rune(value)) != ln {
				return ErrValidateLen
			}
		case "regexp":
			matched, err := regexp.MatchString(rule, value)
			if err != nil {
				return ErrInvalidValidationRule
			}

			if !matched {
				return ErrValidateRegexp
			}
		}
	}

	return nil
}

func validateInt(value int64, rules validationRules) error {
	for key, rule := range rules {
		switch key {
		case "in":
			collection := make([]int64, 0)
			collectionRaw := strings.Split(rule, ",")

			for _, item := range collectionRaw {
				valueRaw, err := strconv.Atoi(item)
				if err != nil {
					return ErrInvalidValidationRule
				}
				collection = append(collection, int64(valueRaw))
			}

			if !slices.Contains(collection, value) {
				return ErrValidateIn
			}
		case "min":
			ln, err := strconv.Atoi(rule)
			if err != nil {
				return ErrInvalidValidationRule
			}

			if value < int64(ln) {
				return ErrValidateMin
			}
		case "max":
			ln, err := strconv.Atoi(rule)
			if err != nil {
				return ErrInvalidValidationRule
			}

			if value > int64(ln) {
				return ErrValidateMax
			}
		}
	}
	return nil
}

func validateSlice(value reflect.Value, rules validationRules) error {
	for i := 0; i < value.Len(); i++ {
		item := value.Index(i)
		switch {
		case item.Kind() == reflect.String:
			err := validateString(item.String(), rules)
			if err != nil {
				return err
			}
		case item.Kind() == reflect.Int:
			err := validateInt(item.Int(), rules)
			if err != nil {
				return err
			}
		default:
			return ErrValueNotSupported
		}
	}

	return nil
}

func extractValidationRules(str string) validationRules {
	rules := make(validationRules)
	rulesRaw := strings.Split(str, "|")

	for _, item := range rulesRaw {
		ruleRaw := strings.Split(item, ":")
		if len(ruleRaw) != 2 {
			continue
		}
		rules[ruleRaw[0]] = ruleRaw[1]
	}

	return rules
}
