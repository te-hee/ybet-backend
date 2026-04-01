package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidatorStruct struct{
	validate *validator.Validate
}


func (v *ValidatorStruct)Validate(out any) error {
    return v.validate.Struct(out)
}

func NewValidateStruct() *ValidatorStruct{
	return &ValidatorStruct{validate: validator.New()}
}

func formatError(errs validator.ValidationErrors) string {
	
	var errorMessages []string

	for _, err := range errs {
		switch err.Tag() {
				case "required":
						errorMessages = append(errorMessages, err.Field() + " is requierd")
				case "uuid":
						errorMessages = append(errorMessages, err.Field() + " needs to be in a uuid format")
				case "min":
						errorMessages = append(errorMessages, err.Field() + " minimal length is "+ err.Param())
				case "max":
						errorMessages = append(errorMessages, err.Field() + " maximal lengtth is " + err.Param())
				case "oneof":
						errorMessages = append(errorMessages, err.Field() + " needs to be of of thoe values: " + err.Param())
				default:
						errorMessages = append(errorMessages, err.Field() + " is invalid")
				}
	}
	
	return  strings.Join(errorMessages, ", ")
}


