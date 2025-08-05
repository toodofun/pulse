package controller

import (
	"github.com/gin-gonic/gin"
	"pulse/internal/model"
	"pulse/internal/service"
)

type MonitorController struct {
	svc *service.MonitorService
}

func NewMonitorController(svc *service.MonitorService) *MonitorController {
	return &MonitorController{
		svc: svc,
	}
}

func (c *MonitorController) handleGetDailyRatioFromPublic(ctx *gin.Context) {
	res, err := c.svc.GetDailySuccessRatios(ctx.Param("id"), GetUser(ctx), true)
	if err != nil {
		Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
		return
	}
	Reply(ctx, CodeSuccess, res)
}

func (c *MonitorController) handleGetDailyRatio(ctx *gin.Context) {
	res, err := c.svc.GetDailySuccessRatios(ctx.Param("id"), GetUser(ctx), false)
	if err != nil {
		Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
		return
	}
	Reply(ctx, CodeSuccess, res)
}

func (c *MonitorController) handleDeleteService(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := c.svc.DeleteService(id, GetUser(ctx)); err != nil {
		Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
		return
	}
	Reply(ctx, CodeSuccess, nil)
}

func (c *MonitorController) handleListServices(ctx *gin.Context) {
	services, err := c.svc.ListServices(GetUser(ctx))
	if err != nil {
		Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
		return
	}
	Reply(ctx, CodeSuccess, services)
}

func (c *MonitorController) handleAddService(ctx *gin.Context) {
	var req model.Service
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Reply(ctx, CodeParamError, nil)
		return
	}
	if err := c.svc.AddService(&req, GetUser(ctx)); err != nil {
		Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
	} else {
		Reply(ctx, CodeSuccess, nil)
	}
}

func (c *MonitorController) handleGetService(ctx *gin.Context) {
	id := ctx.Param("id")
	if res, err := c.svc.GetServiceByID(id, GetUser(ctx)); err != nil {
		Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
	} else {
		Reply(ctx, CodeSuccess, res)
	}
}

func (c *MonitorController) handleSetEnable(ctx *gin.Context) {
	if err := c.svc.SetEnabled(ctx.Param("id"), true, GetUser(ctx)); err != nil {
		Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
	} else {
		Reply(ctx, CodeSuccess, nil)
	}
}

func (c *MonitorController) handleSetDisable(ctx *gin.Context) {
	if err := c.svc.SetEnabled(ctx.Param("id"), false, GetUser(ctx)); err != nil {
		Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
	} else {
		Reply(ctx, CodeSuccess, nil)
	}
}

func (c *MonitorController) handleUpdateService(ctx *gin.Context) {
	var req model.Service
	if err := ctx.ShouldBindJSON(&req); err != nil {
		Reply(ctx, CodeParamError, nil)
		return
	}
	if err := c.svc.UpdateService(&req, ctx.Param("id"), GetUser(ctx)); err != nil {
		Reply(ctx, NewCodeWithMsg(CodeUnknown, err.Error()), nil)
	} else {
		Reply(ctx, CodeSuccess, nil)
	}
}

func (c *MonitorController) RegisterRoute(group *gin.RouterGroup) {
	api := group.Group("/monitor")
	api.GET("/:id/daily", c.handleGetDailyRatio)
	api.GET("/:id/daily/public", c.handleGetDailyRatioFromPublic)
	api.GET("/:id/detail", c.handleGetService)
	api.DELETE("/:id", c.handleDeleteService)
	api.PUT("/:id", c.handleUpdateService)
	api.PUT("/:id/enable", c.handleSetEnable)
	api.PUT("/:id/disable", c.handleSetDisable)
	api.GET("", c.handleListServices)
	api.POST("", c.handleAddService)
}
