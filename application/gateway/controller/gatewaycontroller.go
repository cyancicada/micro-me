package controller

import (
	"github.com/gin-gonic/gin"

	"micro-me/application/common/baseresponse"
	"micro-me/application/gateway/logic"
)

type (
	GateController struct {
		gateLogic *logic.GateWayLogic
	}
)

func NewGateController(gateLogic *logic.GateWayLogic) *GateController {

	return &GateController{gateLogic: gateLogic}
}

func (c *GateController) Send(context *gin.Context) {
	r := new(logic.SendRequest)
	if err := context.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(context, err)
		return
	}
	res, err := c.gateLogic.Send(r)
	baseresponse.HttpResponse(context, res, err)
	return
}


func (c *GateController) GetServerAddress(context *gin.Context) {
	r := new(logic.GetServerAddressRequest)
	if err := context.ShouldBindJSON(r); err != nil {
		baseresponse.ParamError(context, err)
		return
	}
	res, err := c.gateLogic.GetServerAddress(r)
	baseresponse.HttpResponse(context, res, err)
	return
}
