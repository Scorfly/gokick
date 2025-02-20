package gokick

import (
	"context"
	"net/http"
)

type PublicKeyResponseWrapper Response[PublicKeyResponse]

type PublicKeyResponse struct {
	PublicKey string `json:"public_key"`
}

func (c *Client) GetPublicKey(ctx context.Context) (PublicKeyResponseWrapper, error) {
	response, err := makeRequest[PublicKeyResponse](
		ctx,
		c,
		http.MethodGet,
		"/public/v1/public-key",
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return PublicKeyResponseWrapper{}, err
	}

	return PublicKeyResponseWrapper(response), nil
}
