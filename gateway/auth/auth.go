package auth

import (
	"context"
	"encoding/base64"
	"fmt"
)

type BasicAuth struct {
	Username string
	Password string
}

func (b BasicAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	auth := fmt.Sprintf("%s:%s", b.Username, b.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))

	return map[string]string{
		"authorization": "Basic " + encoded,
	}, nil
}

func (b BasicAuth) RequireTransportSecurity() bool {
	return false
}