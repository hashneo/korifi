package presenter

import (
	"code.cloudfoundry.org/korifi/api/config"
)

const V3APIVersion = "3.117.0+cf-k8s"

type APILink struct {
	Link
	Meta *APILinkMeta `json:"meta,omitempty"`
}

type APILinkMeta struct {
	Version string `json:"version,omitempty"`
}

type RootResponse struct {
	Links   map[string]*APILink `json:"links"`
	CFOnK8s bool                `json:"cf_on_k8s"`
}

type RootBuilder struct {
	ApiConfig *config.APIConfig
}

func (b *RootBuilder) ServerRoot() string {
	return b.ApiConfig.ServerURL
}

func (b *RootBuilder) UaaUrl() string {
	return b.ApiConfig.UaaUrl
}

func (b *RootBuilder) LoginUrl() string {
	return b.ApiConfig.LoginUrl
}

func (b *RootBuilder) CFOnK8s() bool {
	return b.ApiConfig.CFOnK8s == nil || *b.ApiConfig.CFOnK8s
}

func (b *RootBuilder) Get() RootResponse {
	rr := RootResponse{
		Links: map[string]*APILink{
			"self":                {Link: Link{HRef: b.ServerRoot()}},
			"bits_service":        nil,
			"cloud_controller_v2": nil,
			"cloud_controller_v3": {Link: Link{HRef: b.ServerRoot() + "/v3"}, Meta: &APILinkMeta{Version: V3APIVersion}},
			"network_policy_v0":   nil,
			"network_policy_v1":   nil,
			"login":               {Link: Link{HRef: b.LoginUrl()}},
			"uaa":                 {Link: Link{HRef: b.UaaUrl()}},
			"credhub":             nil,
			"routing":             nil,
			"logging":             nil,
			"log_cache":           {Link: Link{HRef: b.ServerRoot()}},
			"log_stream":          nil,
			"app_ssh":             nil,
		},
		CFOnK8s: b.CFOnK8s(),
	}

	return rr
}
