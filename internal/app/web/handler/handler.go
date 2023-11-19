package handler

import "github.com/fakovacic/editor/internal/app/web"

type Handler struct {
	service web.Service
}

func New(service web.Service) *Handler {
	return &Handler{
		service: service,
	}
}
