package entity

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"net/url"
	"strconv"
	"strings"
)

type SeatStatus string

const (
	StatusFree     SeatStatus = "FREE"
	StatusReserved SeatStatus = "RESERVED"
	StatusSold     SeatStatus = "SOLD"
)

func (s SeatStatus) Valid() bool {
	switch s {
	case StatusFree, StatusReserved, StatusSold:
		return true
	default:
		return false
	}
}

type SeatsListQuery struct {
	pagination
	EventIDRaw   string `json:"event_id"`
	RowNumberRaw string `json:"row_number"`
	StatusRaw    string `json:"status"`

	// нормализованные значения
	EventID   int64       `json:"-"`
	RowNumber *int        `json:"-"`
	Status    *SeatStatus `json:"-"`
}

func (c *SeatsListQuery) Prepare(values url.Values) error {
	// сырье
	c.EventIDRaw = strings.TrimSpace(values.Get("event_id"))
	c.RowNumberRaw = strings.TrimSpace(values.Get("row_number"))
	c.StatusRaw = strings.TrimSpace(values.Get("status"))

	// нормализация status
	if c.StatusRaw != "" {
		normalized := SeatStatus(strings.ToUpper(c.StatusRaw))
		c.Status = &normalized
	}

	// парсинг чисел (включая soft-empty для row_number)
	if c.EventIDRaw != "" {
		id, err := strconv.ParseInt(c.EventIDRaw, 10, 64)
		if err != nil {
			return fmt.Errorf("event_id: должно быть целым числом")
		}
		c.EventID = id
	}
	if c.RowNumberRaw != "" {
		rn, err := strconv.Atoi(c.RowNumberRaw)
		if err != nil {
			return fmt.Errorf("row_number: должно быть целым числом")
		}
		c.RowNumber = &rn
	}

	// пагинация
	c.pagination.Prepare(values)

	return c.Validate()
}

func (c *SeatsListQuery) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.pagination),
		validation.Field(
			&c.EventID,
			validation.Required.Error("ID события: обязательное поле."),
			validation.Min(int64(1)).Error("ID события должно быть больше нуля"),
		),
		validation.Field(
			&c.RowNumber,
			validation.Min(1).Error("row_number должен быть >= 1"),
			validation.Max(5000).Error("row_number слишком велик"),
		),
	)
}

// GetCacheKey формирует стабильный и короткий ключ кэша.
// Требование: вызывай после Prepare() и Validate(), чтобы поля уже были нормализованы.
func (c *SeatsListQuery) GetCacheKey() string {
	const (
		ns      = "seats:list" // пространство имён/ресурс
		version = "v1"         // меняй при изменении логики формирования ключа
	)

	// Канонизируем опциональные части, чтобы отсутствие и пустая строка давали одинаковый результат.
	rowPart := "row:-" // "-" = отсутствует
	if c.RowNumber != nil {
		rowPart = "row:" + strconv.Itoa(*c.RowNumber)
	}
	statusPart := "status:-"
	if c.Status != nil {
		// гарантируем верхний регистр (на случай ручного создания объекта)
		statusPart = "status:" + strings.ToUpper(string(*c.Status))
	}

	// Собираем компактную строку параметров. Порядок фиксированный.
	// Включаем только значения, влияющие на ответ.
	payload := strings.Join([]string{
		"page:" + strconv.Itoa(c.Page),
		"size:" + strconv.Itoa(c.Limit),
		rowPart,
		statusPart,
	}, "|")

	// Хешируем хвост, чтобы ключи были короткими и не зависящими от длины payload.
	sum := sha1.Sum([]byte(payload))
	// Если хочешь укоротить — возьми первые 12 байт, но полный sha1 безопаснее от коллизий.
	hashHex := hex.EncodeToString(sum[:])

	// Финальный ключ. event_id выводим отдельно — так проще делать точечную инвалидацию по событию.
	// Пример: seats:list:v1:ev:123:9f0c...e2a
	return fmt.Sprintf("%s:%s:ev:%d:%s", ns, version, c.EventID, hashHex)
}

type SeatForList struct {
	ID     int64  `json:"id"`
	Row    int    `json:"row"`
	Number int    `json:"number"`
	Status string `json:"status"`
}
