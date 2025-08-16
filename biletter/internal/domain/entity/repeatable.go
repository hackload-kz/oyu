package entity

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type pagination struct {
	Limit  int
	Offset int
	Page   int
}

func (p *pagination) Prepare(values url.Values) {
	limit, err := strconv.Atoi(values.Get("pageSize"))
	if err != nil {
		limit = 10
		err = nil
	}

	page, err := strconv.Atoi(values.Get("page"))
	if err != nil {
		page = 1
		err = nil
	}

	p.Offset = (page - 1) * limit
	p.Limit = limit
	p.Page = page
}

func (p *pagination) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(
			&p.Limit,
			validation.Max(20).Error("Максимально допустимый лимит записей - 20!"),
		),
	)
}

func (p *pagination) GenerateFilters() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.Limit, p.Offset)
}

func (p *pagination) GenerateFiltersBetween() string {
	var start int

	if p.Page != 1 {
		start = p.Offset + 1
	}

	end := p.Offset + p.Limit

	return fmt.Sprintf("WHERE rn BETWEEN %d AND %d", start, end)
}

// specificValidation валидация спецсимволов во входящих параметрах
func specificValidation(value any) (err error) {
	re, err := regexp.Compile(`(?i)union.*select|select.*from|insert.*into|drop.*table|drop.*database|update.*set|--|\*|\$|\|\||%|\\n|\\r|\\t|\\b|\\z|\\A`)
	if err != nil {
		return err
	}

	var s string
	var ok bool
	if s, ok = value.(string); !ok {
		var sp *string
		if sp, ok = value.(*string); !ok {
			return fmt.Errorf("не удалось опознать значение")
		}

		if sp == nil {
			return nil
		}

		s = *sp
	}

	matched := re.FindAllString(s, -1)

	// проверка амперсандов (не сущностей)
	amp := strings.Count(s, "&")
	validEntityRe := regexp.MustCompile(`&[a-zA-Z]+;|&#\d+;|&#x[0-9a-fA-F]+;`)
	validEntities := validEntityRe.FindAllString(s, -1)

	if amp > len(validEntities) {
		matched = append(matched, "&")
	}

	// проверка одиночного #
	if strings.Contains(s, "#") {
		entityHashRe := regexp.MustCompile(`&#\d+;|&#x[0-9a-fA-F]+;`)
		if !entityHashRe.MatchString(s) {
			matched = append(matched, "#")
		}
	}

	if len(matched) != 0 {
		errForUser := errors.New("убедитесь что вы не ввели запрещённые символы: &,#,*,$,\\,|,%,\\n,\\t,\\r,\\b,\\z,\\A")
		errForLog := fmt.Errorf("%w %s", errForUser, matched)
		return errForLog
	}

	return nil
}

const layoutYMD = "2006-01-02"

// кастомное правило: пусто — ок; иначе строго YYYY-MM-DD
func validISODateYYYYMMDD(v interface{}) error {
	s, _ := v.(string)
	if s == "" {
		return nil
	}
	if _, err := time.Parse(layoutYMD, s); err != nil {
		return errors.New("validation_date date must be in format YYYY-MM-DD")
	}
	return nil
}

func RuneCount(s string) int { return len([]rune(s)) }

// безопасный SQL literal: 'x' -> ”x”
func Quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// экранирование для LIKE: %, _, \ — спецсимволы
func EscapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}
