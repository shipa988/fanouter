package controllers

import (
	"github.com/shipa988/fanouter/internal/domain/usecase/sender"
)

var _ sender.QuerySenderFabric = (*HTTPClientFabric)(nil)

type HTTPClientFabric struct {
}

func NewHTTPClientFabric() *HTTPClientFabric {
	return &HTTPClientFabric{}
}

func (H HTTPClientFabric) NewQuerySender() sender.QuerySender {
	return &HTTPClient{}
}
