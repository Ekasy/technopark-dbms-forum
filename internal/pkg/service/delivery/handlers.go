package delivery

import (
	"forum/internal/models"
	"forum/internal/pkg/service"
	"net/http"

	"github.com/gorilla/mux"
)

type ServiceDelivery struct {
	serviceUsecase service.ServiceUsecase
}

func NewServiceDelivery(serviceUsecase service.ServiceUsecase) *ServiceDelivery {
	return &ServiceDelivery{
		serviceUsecase: serviceUsecase,
	}
}

func (sd *ServiceDelivery) Routing(r *mux.Router) {
	r.HandleFunc("/service/status", sd.GetServiceStatusHandler).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/service/clear", sd.ClearServiceHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (sd *ServiceDelivery) GetServiceStatusHandler(w http.ResponseWriter, r *http.Request) {
	status, err := sd.serviceUsecase.GetServiceStatus()
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
		w.Write(models.ToBytes(status))
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}

func (sd *ServiceDelivery) ClearServiceHandler(w http.ResponseWriter, r *http.Request) {
	err := sd.serviceUsecase.ClearService()
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(models.ToBytes(models.Error{Message: err.Error()}))
	}
}
