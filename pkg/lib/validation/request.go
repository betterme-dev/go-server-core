package validation

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type (
	Errors []Error
	Error  struct {
		Field   string `json:"field,omitempty"`
		Message string `json:"message"`
		Tag     string `json:"-"`
	}
)

func (e *Errors) Add(err Error) {
	*e = append(*e, err)
}

func (e Error) Error() string {
	return e.Message
}

const (
	defaultTag = "binding"
)

func ValidateRequest(req interface{}) error {
	v := validator.New()
	v.SetTagName(defaultTag)
	return v.Struct(req)
}
func ValidateRequestWithTag(req interface{}, tag string) error {
	v := validator.New()
	if tag == "" {
		tag = defaultTag
	}
	v.SetTagName(tag)
	return v.Struct(req)
}

func ErrorsList(err error) Errors {
	var errors = Errors{}
	switch e := err.(type) {
	case *json.UnmarshalTypeError:
		msg := fmt.Sprintf("Field-type of %s expects to be %s", e.Field, e.Type.String())
		errors.Add(Error{Field: e.Field, Message: msg, Tag: "type"})
	case validator.ValidationErrors:
		for _, f := range e {
			message := fmt.Sprintf("Field %s expects to be %s", fieldName(f), f.ActualTag())
			if f.Param() != "" {
				message = fmt.Sprintf("%s with value %s", message, f.Param())
			}
			errors.Add(Error{Field: f.StructField(), Message: message, Tag: f.Tag()})
		}
	case Error:
		errors.Add(e)
	case error:
		errors.Add(Error{Message: e.Error()})
	}
	return errors
}

func fieldName(f validator.FieldError) string {
	name := f.Field()
	if name == "" {
		name = f.StructField()
	}
	return name
}
