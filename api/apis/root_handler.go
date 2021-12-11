package apis

import (
	"net/http"

	"code.cloudfoundry.org/cf-k8s-controllers/api/presenter"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
)

const (
	RootGetEndpoint = "/"
)

type RootHandler struct {
	logger    logr.Logger
	serverURL string
}

func NewRootHandler(logger logr.Logger, serverURL string) *RootHandler {
	return &RootHandler{serverURL: serverURL}
}

func (h *RootHandler) rootGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := writeJsonResponse(w, presenter.GetRootResponse(h.serverURL), http.StatusOK)
	if err != nil {
		h.logger.Error(err, "Failed to render response")
		writeUnknownErrorResponse(w)
	}
}

func (h *RootHandler) RegisterRoutes(router *mux.Router) {
	router.Path(RootGetEndpoint).Methods("GET").HandlerFunc(h.rootGetHandler)
}
