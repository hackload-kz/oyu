package postgresql

import (
	"biletter/internal/domain/entity"
	"context"
	"fmt"
	"strings"
)

// GetSeatsList реализует фильтры: event_id (обяз.), row (опц.), status (опц.), пагинацию.
// Сортировка стабильная: row_number, seat_number, id
func (s *storage) GetSeatsList(ctx context.Context, query entity.SeatsListQuery) ([]entity.SeatForList, error) {
	// safety: event обязателен
	if query.EventID == 0 {
		return nil, fmt.Errorf("event_id is required")
	}

	var (
		args     []any
		cond     []string
		argIndex = 1
	)

	// обязательное условие: event_id
	cond = append(cond, fmt.Sprintf("s.event_id = $%d", argIndex))
	args = append(args, query.EventID)
	argIndex++

	// опциональный фильтр по ряду
	if query.RowNumber != nil {
		cond = append(cond, fmt.Sprintf("s.row_number = $%d", argIndex))
		args = append(args, *query.RowNumber)
		argIndex++
	}

	// опциональный фильтр по статусу
	// status вычисляется на лету:
	// SOLD: sold = true
	// RESERVED: sold = false AND reserved_by_booking IS NOT NULL AND reserved_until > now()
	// FREE: sold = false AND (reserved_by_booking IS NULL OR reserved_until IS NULL OR reserved_until <= now())
	if query.Status != nil {
		switch *query.Status {
		case entity.StatusSold:
			cond = append(cond, "s.sold = true")
		case entity.StatusReserved:
			cond = append(cond,
				"s.sold = false",
				"s.reserved_by_booking IS NOT NULL",
				"s.reserved_until > NOW()",
			)
		case entity.StatusFree:
			cond = append(cond,
				"s.sold = false",
				"(s.reserved_by_booking IS NULL OR s.reserved_until IS NULL OR s.reserved_until <= NOW())",
			)
		default:
			// не должно случиться после Validate, но на всякий
			return nil, fmt.Errorf("unsupported status")
		}
	}

	where := strings.Join(cond, " AND ")

	sql := `
		SELECT
		  	s.id,
		  	s.row_number  AS row,
		  	s.seat_number AS number,
		  	CASE
				WHEN s.sold THEN 'SOLD'
				WHEN s.reserved_by_booking IS NOT NULL AND s.reserved_until > NOW() THEN 'RESERVED'
				ELSE 'FREE'
		  	END AS status
		FROM seats s
		WHERE 
		` + where + `ORDER BY s.row_number, s.seat_number, s.id LIMIT 
		$` + fmt.Sprint(argIndex) + ` OFFSET $` + fmt.Sprint(argIndex+1)

	args = append(args, query.Limit, query.Offset)

	rows, err := s.client.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := make([]entity.SeatForList, 0, query.Limit)
	for rows.Next() {
		var item entity.SeatForList
		if err := rows.Scan(&item.ID, &item.Row, &item.Number, &item.Status); err != nil {
			return nil, err
		}
		response = append(response, item)
	}

	return response, rows.Err()
}
