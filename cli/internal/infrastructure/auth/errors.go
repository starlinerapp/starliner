package auth

import (
	"errors"
	"fmt"
	"net/http"

	openapi "starliner.app/cli/internal/infrastructure/auth/generated/client"
)

func loginError(err error, httpResp *http.Response) error {
	openAPIErr, ok := errors.AsType[*openapi.GenericOpenAPIError](err)
	if !ok || httpResp == nil {
		return err
	}

	switch m := openAPIErr.Model().(type) {
	case openapi.SocialSignIn400Response:
		return fmt.Errorf("status %d: %s", httpResp.StatusCode, m.Message)
	case openapi.SocialSignIn403Response:
		if m.Message != nil {
			return fmt.Errorf("status %d: %s", httpResp.StatusCode, *m.Message)
		}
	}

	if body := openAPIErr.Body(); len(body) > 0 {
		return fmt.Errorf("status %d: %s", httpResp.StatusCode, body)
	}
	return fmt.Errorf("status %d", httpResp.StatusCode)
}
