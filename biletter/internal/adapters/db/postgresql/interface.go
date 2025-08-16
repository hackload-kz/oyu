package postgresql

import (
	"biletter/internal/domain/entity"
	"biletter/internal/grpc_client"
	"biletter/pkg/client/postgresql_client"
	"biletter/pkg/logging"
	"context"
	"errors"
)

type PostgresStorage interface {
	GetUserByEmail(ctx context.Context, email string) (data entity.AuthUser, err error)

	GetEventList(ctx context.Context, filters string) (data []entity.EventForList, err error)

	BookingCreate(ctx context.Context, eventID, userID int64) (bookingID int64, err error)

	GetSeatsList(ctx context.Context, query entity.SeatsListQuery) ([]entity.SeatForList, error)
}

type storage struct {
	client postgresql_client.Client
	logger *logging.Logger
	rpc    grpc_client.GrpcClient
}

func NewPostgresStorage(client postgresql_client.Client, logger *logging.Logger, router grpc_client.GrpcClient) PostgresStorage {
	return &storage{
		client: client,
		logger: logger,
		rpc:    router,
	}
}

func (s *storage) getLang(ctx context.Context) (lang, suffix string, err error) {
	var ok bool

	if lang, ok = ctx.Value("lang").(string); !ok {
		err = errors.New("не удалось получить данные установленного языка")
		s.logger.Error(err)
		return "", "", err
	}

	if lang != "ru" {
		suffix = "_" + lang
	}

	return lang, suffix, nil
}
