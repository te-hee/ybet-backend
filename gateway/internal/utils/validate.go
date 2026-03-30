package utils

import (
	"gateway/internal/model"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

func validateStruct(s interface{}) *model.OutputError {
    var errorsMessages []string
    err := validate.Struct(s)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            errorsMessages = append(errorsMessages, formatError(err))
        }
    }


		if len(errorsMessages) == 0{
			return nil
		}

		errorMessage := strings.Join(errorsMessages, ",")
		error := model.NewOutputError(errorMessage)
    return &error
}

func formatError(err validator.FieldError) string {
    switch err.Tag() {
    case "required":
        return err.Field() + " is requierd"
    case "uuid":
        return err.Field() + " needs to be in a uuid format"
    case "min":
        return err.Field() + " minimal length is "+ err.Param()
    case "max":
        return err.Field() + " maximal lengtth is " + err.Param()
    case "oneof":
				return err.Field() + " needs to be of of thoe values: " + err.Param()
    default:
        return err.Field() + " is invalid"
    }
}

func ValidateQuery[T any](c fiber.Ctx) (*T, *model.OutputError){
	query := new(T)
	if err := c.Bind().Query(&query); err != nil{
		outputErr := model.NewOutputError("Bad query")
		return nil,  &outputErr
	}

	if outputErr := validateStruct(query); outputErr != nil{
		return nil,  outputErr
	}

	return query, nil
}

func ValidateBody[T any](c fiber.Ctx) (*T, *model.OutputError){
	body := new(T)
	if err := c.Bind().Body(&body); err != nil{
		outputErr := model.NewOutputError("Bad json")
		return nil,  &outputErr
	}

	if outputErr := validateStruct(body); outputErr != nil{
		return nil,  outputErr
	}

	return body, nil
}
