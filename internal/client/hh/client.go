package hh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/KOJIMEISTER/it_russian_stat/internal/domain"
	"github.com/KOJIMEISTER/it_russian_stat/pkg/config"
)

var (
	ErrVacancyNotFound = errors.New("vacancy not found")
	ErrRateLimited     = errors.New("rate limited")
)

type HHClient struct {
	Config config.ClientConfig
	Client *http.Client
}

func NewHHClient(cfg config.ClientConfig) *HHClient {
	return &HHClient{
		Config: cfg,
		Client: &http.Client{Timeout: cfg.GetTimeout()},
	}
}

func (hh *HHClient) GetVacancyList(ctx context.Context, reqData *domain.VacancySearchRequest) (*domain.VacancySearchResponse, error) {
	searchQuery := fmt.Sprintf("%s?area=%s&professional_role=%s&date_from=%s&date_to=%s&per_page=%d&page=%d",
		hh.Config.GetSearchUrl(), reqData.Area, reqData.Role, reqData.StartDate, reqData.EndDate, reqData.PerPage, reqData.Page)

	request, err := http.NewRequestWithContext(ctx, "GET", searchQuery, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", hh.Config.GetBearerToken()))

	resp, err := hh.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response domain.VacancySearchResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	ids := make([]string, 0, len(response.Vacancies))
	for _, vac := range response.Vacancies {
		ids = append(ids, vac.ID)
	}

	return &response, nil
}

func (hh *HHClient) GetVacancyDetails(ctx context.Context, vacID string) (map[string]interface{}, error) {
	vacUrl := hh.Config.GetVacancyUrl() + vacID

	req, err := http.NewRequestWithContext(ctx, "GET", vacUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", hh.Config.GetBearerToken()))

	resp, err := hh.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var data map[string]interface{}
		if err = json.Unmarshal(body, &data); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
		return data, nil

	case http.StatusNotFound:
		return nil, fmt.Errorf("vacancy not found: %w", ErrVacancyNotFound)
	case http.StatusForbidden, http.StatusTooManyRequests:
		return nil, fmt.Errorf("rate limited: %w", ErrRateLimited)
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
