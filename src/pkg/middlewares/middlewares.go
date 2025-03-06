package middlewares

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang_lambda_boilerplate/src/pkg/configs"
	"golang_lambda_boilerplate/src/pkg/errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
)

type contextKey string

const (
	// ParamsKey context key to get request body
	ParamsKey contextKey = "params"
)

// HandlerFunc type handler func
type HandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// Generate generate middlewares by list
func Generate(h HandlerFunc, middlewares ...func(HandlerFunc) HandlerFunc) HandlerFunc {
	for i := range middlewares {
		h = middlewares[i](h)
	}

	return h
}

func Begin(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// defer databases.Close()
		defer configs.ClearLF()

		configs.Request(request)

		go func() {
			<-ctx.Done()
			if ctx.Err() == context.DeadlineExceeded {
				configs.LogDeadlineExceeded()
			}
		}()

		return next(ctx, request)
	})
}

func ParseParameters(params interface{}) func(HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return HandlerFunc(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			errorResponse := errors.ErrorResponse{}

			// clear params value before get new value from rerquest
			val := reflect.ValueOf(params).Elem()
			val.Set(reflect.Zero(val.Type()))

			// first: parse path parameters
			// to set path value for field, set field's tag as path, ex: path:"drink_id"
			// required set json tag - (json:"-") to avoid overriding field in next step and hide field when log request body
			err := parse(params, "path", request.PathParameters)
			if err != nil {
				return configs.NotFound(err)
			}

			// second: parse query string parameters
			// to set query param value for field, set field's tag as query, ex: query:"status"
			// required set json tag - (json:"-") to avoid overriding field in next step and hide field when log request body
			err = parse(params, "query", request.QueryStringParameters)
			if err != nil {
				errorResponse.ParseNumberError()
				return configs.BadRequest(errorResponse)
			}

			// third: parse body parameters
			err = parseBody(params, request)
			if err != nil {
				errorResponse.JSONParseError()
				return configs.BadRequest(errorResponse)
			}
			configs.RequestBody(params)

			// validation
			err = configs.Struct(params)
			if err != nil {
				errorResponse.SetByValidationErrors(err.(validator.ValidationErrors))
				return configs.BadRequest(errorResponse)
			}

			newCtx := context.WithValue(ctx, ParamsKey, params)
			return next(newCtx, request)
		})
	}
}

func parseBody(params interface{}, request events.APIGatewayProxyRequest) error {
	if request.Body == "" {
		return nil
	}

	// parse body
	requestBody := request.Body
	if request.IsBase64Encoded {
		decodeString, _ := base64.StdEncoding.DecodeString(requestBody)
		requestBody = string(decodeString)
	}
	return json.Unmarshal([]byte(requestBody), params)
}

func parse(params interface{}, tagName string, m map[string]string) error {
	v := reflect.ValueOf(params).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fmt.Println(field)
		tagValue := field.Tag.Get(tagName)
		fmt.Println(tagValue)
		if tagValue == "" {
			continue
		}
		fmt.Println(v.Field(i))
		err := setFieldValue(v.Field(i), tagValue, m, tagName)
		if err != nil {
			return err
		}
	}
	return nil
}
func setFieldValue(field reflect.Value, tagValue string, m map[string]string, tagName string) error {
	kind := field.Type().Kind()
	fmt.Println(kind)
	if kind == reflect.Ptr {
		kind = field.Type().Elem().Kind()
	}

	if m[tagValue] == "" && tagName == "query" {
		return nil
	}

	switch kind {
	case reflect.Uint64:
		valueUint64, err := strconv.ParseUint(m[tagValue], 10, 32)
		if err != nil {
			return err
		}
		if field.Type().Kind() == reflect.Ptr {
			rv := reflect.ValueOf(&valueUint64)
			field.Set(rv)
		} else {
			field.SetUint(valueUint64)
		}
	case reflect.Int32:
		valueInt64, err := strconv.ParseInt(m[tagValue], 10, 32)
		if err != nil {
			return err
		}
		valueInt32 := int32(valueInt64)
		if field.Type().Kind() == reflect.Ptr {
			rv := reflect.ValueOf(&valueInt32)
			field.Set(rv)
		} else {
			field.SetInt(valueInt64)
		}
	case reflect.Bool:
		field.SetBool(m[tagValue] != "")
	case reflect.Slice:
		if m[tagValue] == "" {
			field.Set(reflect.Zero(field.Type()))
			break
		}

		sliceValue := strings.Split(m[tagValue], ",")
		field.Set(reflect.ValueOf(sliceValue))
	default:
		field.SetString(m[tagValue])
	}
	return nil
}
