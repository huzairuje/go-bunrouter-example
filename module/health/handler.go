package health

import (
	"net/http"

	"go-bunrouter-example/infrastructure/httplib"
	logger "go-bunrouter-example/infrastructure/log"
	"go-bunrouter-example/module/primitive"
	"go-bunrouter-example/utils"

	"github.com/uptrace/bunrouter"
)

type Http struct {
	serviceHealth InterfaceService
}

func NewHttp(serviceHealth InterfaceService) InterfaceHttp {
	return &Http{
		serviceHealth: serviceHealth,
	}
}

type InterfaceHttp interface {
	GroupHealth(group *bunrouter.Group)
}

func (h *Http) GroupHealth(g *bunrouter.Group) {
	g.GET("/ping", h.Ping)
	g.GET("/check", h.HealthCheckApi)
}

func (h *Http) Ping(w http.ResponseWriter, r bunrouter.Request) error {
	return httplib.SetSuccessResponse(w, http.StatusOK, http.StatusText(http.StatusOK), "pong")
}

func (h *Http) HealthCheckApi(w http.ResponseWriter, r bunrouter.Request) error {
	logCtx := "handler.HealthCheckApi"
	ctx := r.Context()
	resp, err := h.serviceHealth.CheckUpTime(ctx)
	if err != nil {
		logger.Error(ctx, utils.ErrorLogFormat, err.Error(), logCtx, "h.serviceHealth.CheckUpTime")
		return httplib.SetErrorResponse(w, http.StatusInternalServerError, primitive.SomethingWentWrong)
	}
	return httplib.SetSuccessResponse(w, http.StatusOK, http.StatusText(http.StatusOK), resp)
}
