package helpers

import (
	"golang_lambda_boilerplate/src/internal/models"

	"github.com/jaswdr/faker"
)

func GenerateUsers(n int) []models.User {
	fake := faker.New()

	var users []models.User
	for i := 0; i < n; i++ {
		user := models.User{
			UserID:   fake.Int(),
			Name:     fake.Person().FirstName(),
			Username: fake.Internet().User(),
			Email:    fake.Person().Contact().Email,
		}
		users = append(users, user)
	}

	return users
}
