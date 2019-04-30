package payment

import "github.com/gin-gonic/gin"

type PaymentModule struct {

}

func Init() (*PaymentModule) {
	return &PaymentModule{}
}

func (module *PaymentModule) RegisterRoutes(router *gin.Engine) {

}
