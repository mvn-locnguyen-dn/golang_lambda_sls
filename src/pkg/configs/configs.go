package configs

import (
	"encoding/json"
	"golang_lambda_boilerplate/src/pkg/constants"
	"golang_lambda_boilerplate/src/pkg/errors"
	"golang_lambda_boilerplate/src/pkg/nulls"
	"golang_lambda_boilerplate/src/pkg/utils"

	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

var lc = newLambdaConfig()

type lambdaConfig struct {
	lf logFormatter
	l  *zerolog.Logger
	v  *validator.Validate
}

func newLambdaConfig() *lambdaConfig {
	// -------------------------------------------------------------
	// *zerolog.Logger
	// -------------------------------------------------------------
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	l := zerolog.New(os.Stderr).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
	// -------------------------------------------------------------
	// *validator.Validate
	// -------------------------------------------------------------
	v := validator.New()
	for tag := range customValidations {
		_ = v.RegisterValidation(tag, customValidations[tag])
	}

	// -------------------------------------------------------------
	// add custom type
	// -------------------------------------------------------------
	v.RegisterCustomTypeFunc(nullTimeTypeFunc, nulls.Time{})
	v.RegisterCustomTypeFunc(nullStringTypeFunc, nulls.String{})
	v.RegisterCustomTypeFunc(nullInt32TypeFunc, nulls.Int32{})

	return &lambdaConfig{
		l: &l,
		v: v,
	}
}

// ClearLF clear value log formatter
func ClearLF() {
	lc.lf = logFormatter{}
}

// Struct log params then validate it
func Struct(params interface{}) error {
	return lc.v.Struct(params)
}

// Request set lf.Request without Body
func Request(request events.APIGatewayProxyRequest) {
	lc.lf.Request = requestLog{
		Path:                  request.Path,
		HTTPMethod:            request.HTTPMethod,
		Headers:               request.Headers,
		QueryStringParameters: request.QueryStringParameters,
	}
}

// RequestBody set lf.Request.Body
func RequestBody(body interface{}) {
	lc.lf.Request.Body = body
}

// SetOthers set lf.Others
func SetOthers(others interface{}) {
	lc.lf.Others = others
}

var logCustomFields = func(e *zerolog.Event) {
	e.Any("request", lc.lf.Request).
		Any("response", lc.lf.Response).
		Any("report", lc.lf.Report).
		Any("others", lc.lf.Others)
}

// LogDeadlineExceeded log when context deadline exceeded
func LogDeadlineExceeded() {
	lc.l.Error().
		Func(logCustomFields).
		Msg("timeout")
}

// BadRequest status code 400 bad request validation
func BadRequest(errorResponse errors.ErrorResponse) (events.APIGatewayProxyResponse, error) {
	defer func() {
		lc.l.Info().
			Func(logCustomFields).
			Msg("validation")
	}()
	return handleResponse(http.StatusBadRequest, errorResponse)
}

// BadRequestSystem status code 400 bad request query database
func BadRequestSystem(message string, err error) (events.APIGatewayProxyResponse, error) {
	errorResponse := errors.ErrorResponse{}
	errorResponse.SetSystemError(message, err)
	lc.lf.Report = err.Error()

	defer func() {
		lc.l.Warn().
			Func(logCustomFields).
			Msg("invalid data")
	}()
	return handleResponse(http.StatusBadRequest, errorResponse)
}

// Unauthorized status code 401
func Unauthorized() (events.APIGatewayProxyResponse, error) {
	defer func() {
		lc.l.Info().
			Func(logCustomFields).
			Msg("unauthorized")
	}()
	return handleResponse(http.StatusUnauthorized, errors.Message{Message: constants.MsgUnauthorized})
}

// Forbidden status code 403
func Forbidden() (events.APIGatewayProxyResponse, error) {
	defer func() {
		lc.l.Info().
			Func(logCustomFields).
			Msg("forbidden")
	}()
	return handleResponse(http.StatusForbidden, errors.Message{Message: constants.MsgForbidden})
}

// NotFound status code 404 not found
func NotFound(err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		lc.lf.Report = err.Error()
	}

	defer func() {
		lc.l.Warn().
			Func(logCustomFields).
			Msg("not found")
	}()

	return handleResponse(http.StatusNotFound, errors.Message{Message: constants.MsgNotFound})
}

// Internal status code 500
func Internal(err error) (events.APIGatewayProxyResponse, error) {
	lc.lf.Report = err.Error()

	defer func() {
		lc.l.Error().
			Func(logCustomFields).
			Msg("internal server error")
	}()
	return handleResponse(http.StatusInternalServerError, errors.Message{Message: constants.MsgInternalServerError})
}

// Success status code 200
func Success(body interface{}) (events.APIGatewayProxyResponse, error) {
	defer func() {
		lc.l.Info().
			Func(logCustomFields).
			Msg("success")
	}()
	return handleResponse(http.StatusOK, body)
}

func handleResponse(statusCode int, body interface{}) (events.APIGatewayProxyResponse, error) {
	marshaled := []byte{}
	if body != nil {
		marshaled, _ = json.Marshal(body)
	}

	lc.lf.Response = responseLog{
		AccessControlAllowOrigin: utils.GetAccessControlAllowOrigin(lc.lf.Request.Headers),
		StatusCode:               statusCode,
		Body:                     body,
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      lc.lf.Response.AccessControlAllowOrigin,
			"X-Frame-Options":                  "deny",
			"Cache-Control":                    "no-store",
			"X-Content-Type-Options":           "nosniff",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: string(marshaled),
	}, nil
}

var (
	nullTimeTypeFunc validator.CustomTypeFunc = func(fl reflect.Value) interface{} {
		value := fl.Interface().(nulls.Time)
		var result *time.Time
		if value.Valid {
			result = &value.Time
		}
		return result
	}

	nullStringTypeFunc validator.CustomTypeFunc = func(fl reflect.Value) interface{} {
		value := fl.Interface().(nulls.String)
		var result *string
		if value.Valid {
			result = &value.String
		}
		return result
	}

	nullInt32TypeFunc validator.CustomTypeFunc = func(fl reflect.Value) interface{} {
		value := fl.Interface().(nulls.Int32)
		var result *int32
		if value.Valid {
			result = &value.Int32
		}
		return result
	}

	customValidations = map[string]validator.Func{}
)
