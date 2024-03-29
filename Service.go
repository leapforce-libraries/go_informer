package informer

import (
	"fmt"
	"net/http"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	apiName string = "Informer"
	apiURL  string = "https://api.informer.eu/v1"
)

type Service struct {
	apiKey       string
	securityCode string
	httpService  *go_http.Service
}

// ServiceError contains error info
//
type ErrorResponse struct {
	Error []string `json:"error"`
}

type ServiceConfig struct {
	ApiKey       string
	SecurityCode string
}

func NewService(config *ServiceConfig) (*Service, *errortools.Error) {
	if config.ApiKey == "" {
		return nil, errortools.ErrorMessage("ApiKey not provided")
	}
	if config.SecurityCode == "" {
		return nil, errortools.ErrorMessage("SecurityCode not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		apiKey:       config.ApiKey,
		securityCode: config.SecurityCode,
		httpService:  httpService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := http.Header{}
	header.Set("ApiKey", service.apiKey)
	header.Set("SecurityCode", service.securityCode)
	(*requestConfig).NonDefaultHeaders = &header

	// add error model
	errorResponse := ErrorResponse{}
	(*requestConfig).ErrorModel = &errorResponse

	request, response, e := service.httpService.HttpRequest(requestConfig)
	if len(errorResponse.Error) > 0 {
		e.SetMessage(strings.Join(errorResponse.Error, "\n"))
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", apiURL, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.apiKey
}

func (service *Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) ApiReset() {
	service.httpService.ResetRequestCount()
}
