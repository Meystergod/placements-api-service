package placements

import (
	"encoding/json"
	"github.com/Meystergod/placements-api-service/internal/apperror"
	"github.com/Meystergod/placements-api-service/internal/config"
	"github.com/Meystergod/placements-api-service/internal/handlers"
	"github.com/Meystergod/placements-api-service/pkg/logging"
	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type handler struct {
	logger     *logging.Logger
	cfg        *config.Config
	httpClient *http.Client
}

func NewHandler(logger *logging.Logger, cfg *config.Config, client *http.Client) handlers.Handler {
	return &handler{
		logger:     logger,
		cfg:        cfg,
		httpClient: client,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/placements/request", apperror.Middleware(h.HandlePlacementRequest))
}

func (h *handler) HandlePlacementRequest(w http.ResponseWriter, req *http.Request) error {
	var placementRequest PlacementRequest

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&placementRequest)
	if err != nil {
		return apperror.ErrorDecode
	}

	v := validator.New()
	err = v.Struct(placementRequest)
	if err != nil {
		return apperror.NewAppError(nil, err.Error(), "", "AS-000100")
	}

	data, err := createPartnerRequestData(h.logger, placementRequest)
	if err != nil {
		return apperror.ErrorEncode
	}

	responsesList := handlePartnerResponse(h.logger, h.httpClient, h.cfg.HTTP.Partners, data)

	response, err := createPlacementResponse(h.logger, placementRequest, responsesList)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)

	return nil
}
