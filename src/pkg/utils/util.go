package utils

import (
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type IUtil interface {
	ConvertStruct(source interface{}, destination interface{}) error
}

type Utils struct{}

func (*Utils) ConvertStruct(source interface{}, destination interface{}) error {
	config := &mapstructure.DecoderConfig{
		Result:           destination,
		TagName:          "json",
		WeaklyTypedInput: true,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	if err := decoder.Decode(source); err != nil {
		return err
	}

	return nil
}

func GetAccessControlAllowOrigin(headers map[string]string) string {
	// set origin = header origin
	// if header has no origin and Referer also, return empty
	// else set origin = header Referer remove backslash at the end
	origin, ok := headers["origin"]
	if !ok {
		referer := headers["Referer"]
		lenReferer := len(referer)
		if lenReferer == 0 {
			return ""
		}
		origin = referer[:lenReferer-1]
	}

	// if ALLOW_ORIGINS is * or origin is include by allowOrigins return origin, else return empty
	if os.Getenv("ALLOW_ORIGINS") == "*" {
		return origin
	}
	allowOrigins := strings.Split(os.Getenv("ALLOW_ORIGINS"), " ")
	for i := range allowOrigins {
		if origin == allowOrigins[i] {
			return origin
		}
	}

	return ""
}
