package handlers

import (
	"code.cloudfoundry.org/korifi/api/config"
	"code.cloudfoundry.org/korifi/api/presenter"
	"net/http"

	"code.cloudfoundry.org/korifi/api/routing"
)

const (
	RootPath = "/"
)

type Root struct {
	apiConfig   *config.APIConfig
	rootBuilder *presenter.RootBuilder
}

func NewRoot(apiConfig *config.APIConfig) *Root {
	return &Root{
		apiConfig:   apiConfig,
		rootBuilder: &presenter.RootBuilder{ApiConfig: apiConfig},
	}
}

func (h *Root) get(r *http.Request) (*routing.Response, error) {
	return routing.NewResponse(http.StatusOK).WithBody(h.rootBuilder.Get()), nil
}

func (h *Root) UnauthenticatedRoutes() []routing.Route {
	return []routing.Route{
		{Method: "GET", Pattern: RootPath, Handler: h.get},
	}
}

func (h *Root) AuthenticatedRoutes() []routing.Route {
	return nil
}
