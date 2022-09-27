package websocket

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type WS_message struct {
	To     string `json:"to" validate:"required"`
	From   string `json:"from" validate:"required"`
	Action string `json:"action" validate:"required,is-action,alpha"`
	Table  string `json:"table" validate:"required,is-table"`
	// Allow for key-value pair for UPDATES or WHERE actions
	Fields map[string]string `json:"fields" validate:""`
}

// Template for inheriting validator checks for actions
func IsAction(fl validator.FieldLevel) bool {
	// Implements validator function for template being supplied
	s := fl.Field().String()
	switch s {
	case "select",
		"update",
		"delete",
		"create":
		return true
	}
	return false
}

// Template for inheriting validator checks for tables
func IsTable(fl validator.FieldLevel) bool {
	// Implements validator function for template being supplied
	s := fl.Field().String()
	switch s {
	case "hosts",
		"network_interfaces",
		"net_flow":
		return true
	}
	return false
}

func NewMessage(json_data []byte) (*WS_message, error) {
	message := new(WS_message)

	validate = validator.New()
	validate.RegisterValidation("is-action", IsAction)
	validate.RegisterValidation("is-table", IsTable)

	json.Unmarshal(json_data, &message)
	err := validate.Struct(message)
	if err != nil {
		return nil, err
	} else {
		return message, nil
	}
}
