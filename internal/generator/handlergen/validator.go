package handlergen

import (
	"errors"
	"fmt"
	"github.com/mrbryside/go-generate/internal/myhttp"
	"github.com/mrbryside/go-generate/internal/mymap"
	"regexp"
	"strings"
)

//TODO: validate type that not nested (map[string]string should not have |required in key) allow only nested*mymap.OrderedMap

// ValidateHandler the handler template data validation function.
func ValidateHandler(handler HandlerTemplateData) error {
	if err := validateType(handler.Type); err != nil {
		return err
	}
	if err := validateName(handler.Name); err != nil {
		return err
	}
	if err := validateApi(handler.Api); err != nil {
		return err
	}
	if err := validateMethod(handler.Method); err != nil {
		return err
	}
	if err := validateRequest(handler.Request); err != nil {
		return err
	}
	if err := validateResponse(handler.Response); err != nil {
		return err
	}
	return nil
}

// Validate the "type" field (should be "handler" or "repository").
func validateType(t string) error {
	if t != "handler" && t != "repository" {
		return errors.New("[type] type must be 'handler' or 'repository'")
	}
	return nil
}

// Validate the "name" field (camel case, starting with lowercase).
func validateName(name string) error {
	if len(name) == 0 || !isCamelCase(name) {
		return errors.New("[name] name must be in camelCase with the first letter lowercase")
	}
	return nil
}

// Validate the "api" field (RESTful style).
func validateApi(api string) error {
	matched, _ := regexp.MatchString(`^/([a-zA-Z0-9_-]+/)*[a-zA-Z0-9_-]+$`, api)
	if !matched {
		return errors.New("[api] api must be in RESTful style, e.g., '/products'")
	}
	return nil
}

// Validate the "method" field (HTTP methods).
func validateMethod(method string) error {
	validMethods := map[string]bool{
		"GET": true, "POST": true, "PUT": true, "DELETE": true, "PATCH": true,
	}
	if !validMethods[strings.ToUpper(method)] {
		return errors.New("[method] method must be a valid HTTP method (GET, POST, PUT, DELETE, PATCH)")
	}
	return nil
}

// Validate the "request" or "response" field with multiple key-value pairs.
func validateStructRequestOrResponse(structType string, content *mymap.OrderedMap) error {
	// Regular expression to validate keys (must start with a letter and contain only letters or numbers with snake_case).
	keyRegex := regexp.MustCompile(`^[a-z]+(_?[a-z0-9]+)*(\|.*)?$`)
	// Regular expression to validate values (either "[gotype]" or "[gotype]|[any string]").
	valueRegex := regexp.MustCompile(`^(bool|byte|complex64|complex128|error|float32|float64|int|int8|int16|int32|int64|rune|string|uint|uint8|uint16|uint32|uint64|uintptr)(\|[^\s\|]+)?$`)

	iter := content.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		key := pair.Key
		value := pair.Value
		// Check if the key matches the required format.
		if !keyRegex.MatchString(key) {
			return fmt.Errorf("[json-spec-format-error][%s] key '%s' is invalid: must start with a letter and contain only letters and numbers or using http status code(200, 201, ...)", structType, key)
		}

		// Check if the value is a nested map and validate it recursively.
		valNested, ok := value.(*mymap.OrderedMap)
		if ok {
			if err := validateStructRequestOrResponse(structType, valNested); err != nil {
				return err
			}
			return nil
		}

		// Check if the value is a slice of nested maps and validate them recursively.
		valSlice, ok := value.([]interface{})
		if ok {
			for _, vs := range valSlice {
				valNestedSlice, ok := vs.(*mymap.OrderedMap)
				if ok {
					if err := validateStructRequestOrResponse(structType, valNestedSlice); err != nil {
						return err
					}
					return nil
				}
				valStringSlice, ok := vs.(string)
				if ok && !valueRegex.MatchString(valStringSlice) {
					return fmt.Errorf("[json-spec-format-error][%s] value for key '%s' is invalid: must be in the format '[gotype]' or '[gotype]|[go validation string]'", structType, key)
				}
				if !ok {
					return fmt.Errorf("[json-spec-format-error][%s] value for key '%s' is invalid: must be string in the format '[gotype]' or '[gotype]|[go validation string]'", structType, key)
				}
			}
			return nil
		}

		// Check if the value is a string and matches the required format.
		valString, ok := value.(string)
		if ok && !valueRegex.MatchString(valString) {
			return fmt.Errorf("[json-spec-format-error][%s] value for key '%s' is invalid: must be in the format '[gotype]' or '[gotype]|[go validation string]'", structType, key)
		}
		if !ok {
			return fmt.Errorf("[json-spec-format-error][%s] value for key '%s' is invalid: must be string in the format '[gotype]' or '[gotype]|[go validation string]'", structType, key)
		}
	}
	return nil
}

func validateRequest(request *mymap.OrderedMap) error {
	err := validateStructRequestOrResponse("request", request)
	if err != nil {
		return err
	}

	return nil
}

// Validate the "response" field, allowing HTTP status codes or custom request-style keys.
func validateResponse(response *mymap.OrderedMap) error {
	if isStatusCodeStyle(response) {
		iter := response.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			key := pair.Key
			if len(key) != 3 {
				return errors.New("[json-spec-format-error][response] HTTP status code must be 3 digits")
			}
			if !myhttp.IsKeyInStatusMapping(key) {
				return errors.New("[json-spec-format-error][response] HTTP status code must not be in the status mapping")
			}

		}

		iter = response.EntriesIter()
		for {
			pair, ok := iter()
			if !ok {
				break
			}
			value := pair.Value
			valNested, ok := value.(*mymap.OrderedMap)
			if ok {
				err := validateStructRequestOrResponse("response", valNested)
				if err != nil {
					return err
				}
			} else {
				return errors.New("[json-spec-format-error][response] value must be an object represent your response code and response body")
			}
		}
	} else {
		err := validateStructRequestOrResponse("response", response)
		if err != nil {
			return err
		}
	}

	return nil
}

// Check if a string is camelCase with the first letter lowercase.
func isCamelCase(s string) bool {
	return regexp.MustCompile(`^[a-z]+[A-Za-z0-9]*$`).MatchString(s)
}

func isStatusCodeStyle(response *mymap.OrderedMap) bool {
	statusCodeRegex := regexp.MustCompile(`^\d{3}$`)
	isStatusCodeStyle := true

	iter := response.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		key := pair.Key
		// Check if the key is a valid 3-digit HTTP status code.
		if !statusCodeRegex.MatchString(key) {
			isStatusCodeStyle = false
		}
	}
	return isStatusCodeStyle
}
