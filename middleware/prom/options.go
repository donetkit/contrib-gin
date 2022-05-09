package prom

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config defines the config for logger middleware
type config struct {
	handlerUrl             string
	namespace              string
	excludeRegexStatus     string
	excludeRegexEndpoint   string
	excludeRegexMethod     string
	endpointLabelMappingFn RequestLabelMappingFn
}

// Option for queue system
type Option func(*config)

// WithNamespace set namespace function
func WithNamespace(namespace string) Option {
	return func(cfg *config) {
		cfg.namespace = namespace
	}
}

// WithHandlerUrl set handlerUrl function
func WithHandlerUrl(handlerUrl string) Option {
	return func(cfg *config) {
		cfg.handlerUrl = handlerUrl
	}
}

// WithExcludeRegexStatus set excludeRegexStatus function
func WithExcludeRegexStatus(excludeRegexStatus string) Option {
	return func(cfg *config) {
		cfg.excludeRegexStatus = excludeRegexStatus
	}
}

// WithExcludeRegexEndpoint set excludeRegexEndpoint function
func WithExcludeRegexEndpoint(excludeRegexEndpoint string) Option {
	return func(cfg *config) {
		cfg.excludeRegexEndpoint = excludeRegexEndpoint
	}
}

// WithExcludeRegexMethod set excludeRegexMethod function
func WithExcludeRegexMethod(excludeRegexMethod string) Option {
	return func(cfg *config) {
		cfg.excludeRegexMethod = excludeRegexMethod
	}
}

// WithEndpointLabelMappingFn set endpointLabelMappingFn function
func WithEndpointLabelMappingFn(endpointLabelMappingFn RequestLabelMappingFn) Option {
	return func(cfg *config) {
		cfg.endpointLabelMappingFn = endpointLabelMappingFn
	}
}

// WithPromHandler set router function
func WithPromHandler(router *gin.Engine) Option {
	return func(cfg *config) {
		if router != nil {
			router.GET(cfg.handlerUrl, promHandler(promhttp.Handler()))
		}
	}
}
