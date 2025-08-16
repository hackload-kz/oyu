package entity

import validation "github.com/go-ozzo/ozzo-validation"

type BookingCreate struct {
	EventID int64 `json:"event_id"`
}

func (c *BookingCreate) Validate() error {
	err := validation.ValidateStruct(c,
		validation.Field(
			&c.EventID,
			validation.Required.Error("ID события: обязательное поле."),
			validation.Min(1).Error("ID события должно быть больше нуля"),
		),
	)

	return err
}
