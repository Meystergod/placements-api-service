package placements

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/Meystergod/placements-api-service/internal/apperror"
	"github.com/Meystergod/placements-api-service/internal/config"
	"github.com/Meystergod/placements-api-service/internal/handlers"
	"github.com/Meystergod/placements-api-service/internal/models/partner"
	"github.com/Meystergod/placements-api-service/internal/models/placement"
	"github.com/Meystergod/placements-api-service/pkg/logging"
	"github.com/Meystergod/placements-api-service/pkg/services"

	"github.com/go-playground/validator"
	"github.com/julienschmidt/httprouter"
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
	router.HandlerFunc(http.MethodPost, config.PLACEMENT_ENDPOINT, apperror.Middleware(h.HandlePlacementRequest))
}

func (h *handler) HandlePlacementRequest(w http.ResponseWriter, req *http.Request) error {
	var placementRequest placement.Request

	h.logger.Info("creating new decoder for placement request body")
	decoder := json.NewDecoder(req.Body)

	h.logger.Info("request body decoding started")
	err := decoder.Decode(&placementRequest)
	if err != nil {
		return apperror.ErrorDecode
	}
	h.logger.Info("request body decoding completed successfully")

	h.logger.Info("request validating started")
	v := validator.New()
	err = v.Struct(placementRequest)
	if err != nil {
		return apperror.ErrorValidate
	}
	h.logger.Info("request validating completed successfully")

	h.logger.Info("creating partner request")
	partnerRequest := placementRequest.ToPartnerRequest()

	h.logger.Info("start getting ads from all partners")
	responsesList := postToAllPartners(h.logger, h.httpClient, partnerRequest, h.cfg.HTTP.Partners)

	h.logger.Info("creating placement response")
	placementResponse := placement.NewResponse(h.logger, placementRequest, responsesList)

	h.logger.Info("start json encoding of placement response")
	data, err := json.Marshal(placementResponse)
	if err != nil {
		h.logger.WithError(err).Error("json encoding of placement response failed")
		return apperror.ErrorEncode
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)

	return nil
}

func postToAllPartners(logger *logging.Logger, client *http.Client, req partner.Request, partners []string) []string {
	c := make(chan string, len(partners))

	var wg sync.WaitGroup

	logger.Info("start sending post request to partners")
	for i := 0; i < len(partners); i++ {
		wg.Add(1)
		go postToPartner(logger, &wg, client, req, c, partners[i])
	}

	wg.Wait()

	close(c)

	logger.Info("converting chan to slice of strings")
	return services.ChanToSlice(c).([]string)
}

func postToPartner(l *logging.Logger, wg *sync.WaitGroup, client *http.Client, r partner.Request, c chan string, partnerUrl string) {
	defer wg.Done()

	url := config.HTTP_URL + partnerUrl + config.PARTNER_ENDPOINT

	logger := l.WithFields(map[string]interface{}{
		"partner url": partnerUrl,
	})

	logger.Info("start json encoding of partner request")
	partnerJSONRequest, err := json.Marshal(r)
	if err != nil {
		logger.WithError(err).Error("json encoding of partner request failed")
		return
	}
	logger.Info("json encoding of partner request completed successfully")

	logger.Info("creating new http post request using partner request json")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(partnerJSONRequest))
	if err != nil {
		logger.WithError(err).Error("creating new http post request failed")
		return
	}
	logger.Info("creating new http post request completed successfully")

	logger.Info("sending new created http post request")
	res, err := client.Do(req)
	if err != nil {
		logger.WithError(err).Error("sending new http post request failed")
		return
	}
	logger.Info("sending new created http post request completed successfully")

	defer res.Body.Close()

	logger.Infof("response status code is %d", res.StatusCode)
	if res.StatusCode == http.StatusNoContent {
		return
	}

	logger.Info("reading all from response body")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.WithError(err).Error("reading all from response body failed")
		return
	}
	logger.Info("reading all from response body completed successfully")

	c <- string(body)
}
