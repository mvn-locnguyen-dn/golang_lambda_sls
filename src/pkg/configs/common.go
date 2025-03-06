package configs

type logFormatter struct {
	Request  requestLog
	Response responseLog
	Report   string
	Others   interface{}
}

type requestLog struct {
	Path                  string            `json:"path"`
	HTTPMethod            string            `json:"httpMethod"`
	Headers               map[string]string `json:"headers"`
	QueryStringParameters map[string]string `json:"queryStringParameters,omitempty"`
	Body                  interface{}       `json:"body,omitempty"`
}

type responseLog struct {
	AccessControlAllowOrigin string      `json:"access_control_allow_origin"`
	StatusCode               int         `json:"statusCode"`
	Body                     interface{} `json:"body,omitempty"`
}
