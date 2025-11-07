package domain

import "time"

type UpdateRequestMessage struct {
	RequestID     string `json:"request_id"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	CallbackQueue string `json:"callback_queue"`
}

type ScrapperResponse struct {
	RequestID  string    `json:"request_id"`
	Status     string    `json:"status"`
	StatusText string    `json:"status_text"`
	At         time.Time `json:"at"`
}

type VacancySearchRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Area      string `json:"area"`
	Role      string `json:"role"`
	Page      int    `json:"page"`
	PerPage   int    `json:"per_page"`
}

type VacancySearchResponse struct {
	Pages     int `json:"pages"`
	Vacancies []struct {
		ID string `json:"id"`
	} `json:"items"`
}
