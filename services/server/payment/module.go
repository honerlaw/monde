package payment

import (
	"github.com/gin-gonic/gin"
	"services/server/payment/service"
	service2 "services/server/core/service"
)

type PaymentModule struct {
	paymentService   *service.PaymentService
}

func Init(countryService *service2.CountryService) (*PaymentModule) {
	paymentService := service.NewPaymentService(countryService)

	return &PaymentModule{
		paymentService:   paymentService,
	}
}

func (module *PaymentModule) RegisterRoutes(router *gin.Engine) {

}
