package middleware

import (
	"encoding/json"
)

var (
	ErrNotFound      = NewAppError(nil, "not found", "", "TS-000004")
	Unauthorized     = NewAppError(nil, "unauthorized", "", "TS-000001")
	FailedDependency = NewAppError(nil, "failed dependency", "", "TS-000024")
	RequestTimeout   = NewAppError(nil, "timeout", "", "DS-000008")
	BadRequest       = NewAppError(nil, "bad request", "", "TS-000400")
	Validation       = NewAppError(nil, "", "", "TS-100400")
	Forbidden        = NewAppError(nil, "forbidden", "", "TS-000003")
)

type AppErr struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

func (e *AppErr) Error() string {
	return e.Message
}

func (e *AppErr) Unwrap() error {
	return e.Err
}

func (e *AppErr) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, developerMessage, code string) *AppErr {
	return &AppErr{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
		Code:             code,
	}
}

func systemError(err error) *AppErr {
	return NewAppError(err, "internal_backup system error", err.Error(), "NS-000000")
}

func FatalUserError(text string, err error) *AppErr {
	return NewAppError(err, text, "", "DS-500010")
}

func PartialUserError(text string, err error) *AppErr {
	return NewAppError(err, text, "", "DS-000002")
}
