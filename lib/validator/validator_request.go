package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

}

func ValidateStruct(s interface{}) error {
	var errorMessage []string
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "email":
				errorMessage = append(errorMessage, "Invalid email format")
			case "required":
				errorMessage = append(errorMessage, "Field "+err.Field()+" wajib diisi")
			case "min":
				if err.Field() == "password" {
					errorMessage = append(errorMessage, "Password minimal 8 character")
				}
			default:
				errorMessage = append(errorMessage, "Field "+err.Field()+" tidak valid")
			}
		}

		return errors.New("Validasi gagal: " + joinMessage(errorMessage))
	}

	return nil
}

func joinMessage(messages []string) string {
	result := ""
	for i, message := range messages {
		if i > 0 {
			result += ", "
		}
		result += message
	}

	return result
}
