package payment

import (
	"github.com/gin-gonic/gin"
	"services/server/payment/service"
)

type PaymentModule struct {
	paymentService   *service.PaymentService
}

func Init() (*PaymentModule) {
	paymentService := service.NewPaymentService()

	return &PaymentModule{
		paymentService:   paymentService,
	}
}

func (module *PaymentModule) RegisterRoutes(router *gin.Engine) {

}
