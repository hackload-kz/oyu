package rest_controller

import (
	"biletter/internal/adapters/db/postgresql"
	"biletter/internal/grpc_client"
	"biletter/internal/middleware"
	"biletter/internal/services"
	"biletter/pkg/logging"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type DocHandler interface {
	Register(*httprouter.Router)

	GetEventList(http.ResponseWriter, *http.Request) error
	BookingCreate(w http.ResponseWriter, r *http.Request) (err error)
	GetSeatsList(w http.ResponseWriter, r *http.Request) (err error)
}

type docHandler struct {
	logger          *logging.Logger
	postgresStorage postgresql.PostgresStorage
	rpc             grpc_client.GrpcClient
	service         services.Service
}

func NewRouterHandler(
	postgresStorage postgresql.PostgresStorage, logger *logging.Logger,
	router grpc_client.GrpcClient, service services.Service,
) DocHandler {
	return &docHandler{
		logger:          logger,
		postgresStorage: postgresStorage,
		rpc:             router,
		service:         service,
	}
}

const (
	eventListURL     = "/api/events"
	bookingCreateURL = "/api/bookings"
	seatsListURL     = "/api/seats"

	swaggerURL = "/api/v1/integrator/swagger/*any"
)

func (h *docHandler) Register(router *httprouter.Router) {
	router.Handler(http.MethodGet, eventListURL, middleware.New(
		middleware.ErrorMiddleware,
	).Then(h.GetEventList, h.postgresStorage, h.logger, "integrator", "r"))

	router.Handler(http.MethodPost, bookingCreateURL, middleware.New(
		middleware.ErrorMiddleware,
		middleware.BasicAuthMiddleware,
	).Then(h.BookingCreate, h.postgresStorage, h.logger, "integrator", "c"))

	router.Handler(http.MethodGet, seatsListURL, middleware.New(
		middleware.ErrorMiddleware,
	).Then(h.GetSeatsList, h.postgresStorage, h.logger, "integrator", "r"))

	router.Handler(http.MethodGet, swaggerURL, httpSwagger.WrapHandler)
}
