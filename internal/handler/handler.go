// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

package handler

import (
	"github.com/tharsis/dashboard-backend/internal/handler/v2"
)

type Handler struct {
	v2 *v2.Handler
}

func New() (*Handler, error) {
	v2Handler, err := v2.NewHandler()
	if err != nil {
		return nil, err
	}

	return &Handler{
		v2: v2Handler,
	}, nil
}
