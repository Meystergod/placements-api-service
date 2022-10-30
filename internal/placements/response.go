package placements

import (
	"bytes"
	"encoding/json"
	"github.com/Meystergod/placements-api-service/internal/apperror"
	"github.com/Meystergod/placements-api-service/internal/config"
	"github.com/Meystergod/placements-api-service/pkg/logging"
	"io"
	"math"
	"net/http"
	"sync"
)

func handlePartnerResponse(logger *logging.Logger, client *http.Client, partners []string, response []byte) []string {
	var partnersLength = len(partners)

	chData := make(chan string, partnersLength)

	var wg sync.WaitGroup

	logger.Info("start goroutines (start post requests to partners)")
	for i := 0; i < partnersLength; i++ {
		wg.Add(1)
		go postPartnerRequest(logger, client, &wg, chData, partners[i], response)
	}
	wg.Wait()

	close(chData)

	return chanToSlice(chData).([]string)
}

func postPartnerRequest(
	l *logging.Logger,
	client *http.Client,
	wg *sync.WaitGroup,
	ch chan string,
	partner string,
	res []byte,
) {
	defer wg.Done()

	url := config.HTTP_URL + partner + config.ENDPOINT

	logger := l.WithFields(map[string]interface{}{
		"PARTNER": partner,
	})

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(res))
	if err != nil {
		logger.WithError(err).Error(apperror.ErrorNewRequestWrap)
		return
	}

	response, err := client.Do(req)
	if err != nil {
		logger.WithError(err).Error(apperror.ErrorSendRequest)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		logger.Info("no content")
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.WithError(err).Error(apperror.ErrorParseBody)
		return
	}

	ch <- string(body)
}

func createPartnerRequestData(logger *logging.Logger, placement PlacementRequest) ([]byte, error) {
	var data PartnerRequest
	dataImp := make([]PartnerRequestImp, len(placement.Tiles))

	for i := 0; i < len(placement.Tiles); i++ {
		dataImp[i].ID = placement.Tiles[i].ID
		dataImp[i].MinWidth = placement.Tiles[i].Width
		dataImp[i].MinHeight = uint(math.Floor(float64(placement.Tiles[i].Width) * placement.Tiles[i].Ratio))
	}

	data.ID = placement.ID
	data.Imp = dataImp
	data.Context = PartnerContext(placement.Context)

	dataJson, err := json.Marshal(data)
	if err != nil {
		logger.Fatal("partner request marshalling error")
		return []byte(""), err
	}

	logger.Info("partner request created")

	return dataJson, nil
}

func createPlacementResponse(logger *logging.Logger, p PlacementRequest, responses []string) ([]byte, error) {
	partnerResponsesList := make([]PartnerResponse, len(responses))
	partnerResponsesSet := make(map[uint][]PartnerResponseImp)

	var wg sync.WaitGroup

	logger.Info("start goroutines to parse partner responses list to map look like id:[imp1, imp2, ...]")
	for i := range responses {
		logger.Info("try to unmarshall partner response", i, "from string to struct")
		if err := json.Unmarshal([]byte(responses[i]), &partnerResponsesList[i]); err != nil {
			logger.WithError(err).Error(apperror.ErrorDecode)
			continue
		}
		wg.Add(1)
		go partnerResponsesToMap(&wg, partnerResponsesList[i], partnerResponsesSet)
	}

	wg.Wait()

	var placementResponse PlacementResponse
	placementResponse.ID = p.ID

	logger.Info("create placement response")
	for j := range p.Tiles {
		if len(partnerResponsesSet[p.Tiles[j].ID]) != 0 {
			partnerMaxValue := getMaxPrice(partnerResponsesSet[p.Tiles[j].ID])
			imp := PlacementResponseImp{
				ID:     partnerMaxValue.ID,
				Width:  partnerMaxValue.Width,
				Height: partnerMaxValue.Height,
				Title:  partnerMaxValue.Title,
				Url:    partnerMaxValue.Url,
			}
			placementResponse.Imp = append(placementResponse.Imp, imp)
		}
	}

	data, err := json.Marshal(placementResponse)
	if err != nil {
		logger.WithError(err).Error(apperror.ErrorEncode)
		return []byte(""), apperror.ErrorEncode
	}

	return data, nil
}
