package http

import (
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan/plugins/imtoken/service"
	"net/http"
)

var (
	svc *service.Service
)

func Router(s *service.Service) []router.Http {
	svc = s
	return []router.Http{
		{"test", system},
	}
}

func system(w http.ResponseWriter, r *http.Request) error {
	return nil
}
