package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

}

// ValidateStruct performs validation on the provided struct using predefined validation rules.
// It checks for various validation tags such as email format, required fields, and minimum length requirements.
//
// Parameters:
//   - s: An interface{} representing the struct to be validated.
//
// Returns:
//   - An error if validation fails, containing a concatenated string of all validation error messages.
//   - nil if the struct passes all validation checks.
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
			case "eqfield":
				errorMessage = append(errorMessage, err.Field()+" harus sama dengan"+err.Param()+".")
			default:
				errorMessage = append(errorMessage, "Field "+err.Field()+" tidak valid")
			}
		}

		return errors.New("Validasi gagal: " + joinMessage(errorMessage))
	}

	return nil
}

// joinMessage concatenates a slice of strings into a single string, separating each message with a comma and space.
//
// Parameters:
//   - messages: A slice of strings containing the messages to be joined.
//
// Returns:
//
//	A string containing all messages from the input slice, joined together with ", " as the separator.
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
