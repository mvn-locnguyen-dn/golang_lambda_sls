package errors

import (
	"fmt"
	"golang_lambda_boilerplate/src/pkg/constants"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Message message response
type Message struct {
	Message string `json:"message"`
}

// ErrorResponse error response
type ErrorResponse struct {
	// example: 店舗コードが必須です
	Message string `json:"message"`
	// example: {"ShopID": "required"}
	Errs map[string]interface{} `json:"errs"`
}

// HasError error response has error
func (er *ErrorResponse) HasError() bool {
	return len(er.Errs) > 0
}

// SetSystemError set system error
func (er *ErrorResponse) SetSystemError(message string, err error) {
	er.Message = message
	er.Errs = map[string]interface{}{
		"system": err.Error(),
	}
}

// SetByValidationErrors set error response by validation errors
func (er *ErrorResponse) SetByValidationErrors(err validator.ValidationErrors) {
	er.Errs = map[string]interface{}{}
	for _, vErr := range err {
		er.Errs[vErr.Field()] = vErr.Tag()
	}
	er.Message = setByFieldError(err[0])
}

func setByFieldError(fieldError validator.FieldError) string {
	tagKey := fieldError.Tag()
	tagMap := map[string]string{}

	tagValue, ok := tagMap[tagKey]
	if !ok {
		panic("not found tag key")
	}
	if !strings.Contains(tagValue, "%s") {
		return tagValue
	}

	// remove [x] because field[x] have the same message as field
	removeChildSymbol := regexp.MustCompile(`\[\d+\]$`)
	fieldKey := removeChildSymbol.ReplaceAllString(fieldError.Field(), "")
	fieldMap := map[string]string{}
	fieldValue, ok := fieldMap[fieldKey]
	if !ok {
		fieldValue = fieldKey
	}
	return fmt.Sprintf(tagValue, fieldValue)
}

func (er *ErrorResponse) ParseNumberError() {
	er.Message = constants.MsgParseNumberError
	er.Errs = map[string]interface{}{
		"number": "number parse error",
	}
}

// JSONParseError set json parse error
func (er *ErrorResponse) JSONParseError() {
	er.Message = constants.MsgJSONParseError
	er.Errs = map[string]interface{}{
		"json": "json parse error",
	}
}
