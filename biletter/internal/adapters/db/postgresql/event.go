package postgresql

import (
	"biletter/internal/domain/entity"
	"context"
	"fmt"
)

func (s *storage) GetEventList(ctx context.Context, filters string) (data []entity.EventForList, err error) {
	query := fmt.Sprintf(`
		SELECT
			id, title
		FROM events
		%s
	`, filters)

	rows, err := s.client.Query(ctx, query)
	if err != nil {
		s.logger.Error(err)
		return data, err
	}

	for rows.Next() {
		var item entity.EventForList

		err = rows.Scan(
			&item.ID,
			&item.Title,
		)

		if err != nil {
			s.logger.Error(err)
			return data, err
		}

		data = append(data, item)
	}

	return data, nil
}
