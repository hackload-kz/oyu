package entity

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"net/url"
	"strings"
	"time"
)

var tzAsiaAlmaty, _ = time.LoadLocation("Asia/Almaty")

type EventListQuery struct {
	Query string `json:"query"`
	Date  string `json:"date"`
}

func (c *EventListQuery) Prepare(values url.Values) (err error) {
	c.Query = values.Get("query")
	c.Date = values.Get("date")

	err = c.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (c *EventListQuery) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Query, validation.By(specificValidation)),
		validation.Field(&c.Date, validation.By(validISODateYYYYMMDD)),
	)
}

type EventForList struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

func (c EventListQuery) GenerateFilters() string {
	var b strings.Builder
	var conds []string

	// 1) Фильтр по дате — диапазон по локальному дню (чётко бьётся в btree по datetime_start)
	if strings.TrimSpace(c.Date) != "" {
		if d, err := time.ParseInLocation(layoutYMD, c.Date, tzAsiaAlmaty); err == nil {
			start := time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, tzAsiaAlmaty)
			end := start.Add(24 * time.Hour)
			conds = append(conds,
				fmt.Sprintf(
					"datetime_start >= TIMESTAMPTZ '%s' AND datetime_start < TIMESTAMPTZ '%s'",
					start.Format(time.RFC3339), end.Format(time.RFC3339),
				),
			)
		}
	}

	// 2) Поиск
	q := strings.TrimSpace(c.Query)
	orderAdded := false
	if q != "" {
		if RuneCount(q) < 3 {
			// короткий запрос — триграммы (idx_events_title_trgm)
			like := "%" + EscapeLike(q) + "%"
			conds = append(conds, "title ILIKE "+Quote(like)+" ESCAPE '\\'")
			b.WriteString(" ORDER BY datetime_start ASC NULLS LAST")
			orderAdded = true
		} else {
			// полнотекстовый — GIN по tsv (idx_events_tsv_gin)
			qLit := Quote(q) // '…' (иммутабельный unaccent используем на стороне БД)
			conds = append(conds,
				"tsv @@ plainto_tsquery('simple', immutable_unaccent("+qLit+"))",
			)
			b.WriteString(" ORDER BY ts_rank(tsv, plainto_tsquery('simple', immutable_unaccent(" + qLit + "))) DESC, datetime_start ASC NULLS LAST")
			orderAdded = true
		}
	}

	// 3) WHERE
	if len(conds) > 0 {
		b.WriteString(" WHERE ")
		b.WriteString(strings.Join(conds, " AND "))
	}

	// 4) ORDER BY по умолчанию
	if !orderAdded {
		b.WriteString(" ORDER BY datetime_start ASC NULLS LAST")
	}

	return b.String()
}
