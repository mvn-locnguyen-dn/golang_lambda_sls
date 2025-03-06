package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang_lambda_boilerplate/src/internal/forms"
	"golang_lambda_boilerplate/src/internal/usecases/users"
	"golang_lambda_boilerplate/src/pkg/configs"
	"golang_lambda_boilerplate/src/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handlerFunc(ctx context.Context, request events.APIGatewayProxyRequest, us users.IUser) (events.APIGatewayProxyResponse, error) {
	var body forms.CreateUserRequest
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		return configs.BadRequestSystem("invalid request body", err)
	}

	fmt.Println(body.UserID)

	response, err := us.Detail()

	switch err {
	case nil:
		return configs.Success(response)
	default:
		return configs.Internal(err)
	}
}

func main() {
	handler := func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return handlerFunc(ctx, request, &users.User{
			Utils: &utils.Utils{},
		})
	}

	lambda.Start(handler)
}
