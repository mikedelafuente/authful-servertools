package httptools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/mikedelafuente/authful-servertools/pkg/config"
	"github.com/mikedelafuente/authful-servertools/pkg/customclaims"
	"github.com/mikedelafuente/authful-servertools/pkg/customerrors"
	"github.com/mikedelafuente/authful-servertools/pkg/logger"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	logger.Debug(ctx, fmt.Sprintf("%s %s", method, url))
	return http.NewRequest(method, url, body)
}

// Does a POST to the specified endpoint. Returns the body bytes, an http status code (0 if no call was made)
func Post(ctx context.Context, url string, requestModel interface{}) ([]byte, int, error) {
	client := &http.Client{}

	logger.Debug(ctx, fmt.Sprintf("Marshaling model for POST %s", url))
	requestBytes, err := MarshalFormat(ctx, requestModel)
	if err != nil {
		logger.Error(ctx, err)
		return nil, 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBytes))
	if err != nil {
		logger.Error(ctx, err)
		return nil, 0, err
	}

	setRequestHeader(ctx, req)

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(ctx, err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(ctx, err)
		return nil, resp.StatusCode, err
	}

	if !IsOkResponse(resp) {
		// TODO: try to extract out an error from the body
		errorMessage := ExtractErrorMessageFromJsonBytes(ctx, bodyBytes, fmt.Sprintf("HTTP Exception calling POST %s returned status %s", url, resp.Status))

		return nil, resp.StatusCode, customerrors.NewServiceError(resp.StatusCode, errorMessage)
	}

	if len(bodyBytes) > 0 {
		logger.Debug(ctx, fmt.Sprintf("Response: %s", string(bodyBytes)))
	}

	return bodyBytes, resp.StatusCode, nil
}

func setRequestHeader(ctx context.Context, req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	jwt := ctx.Value(customclaims.ContextJwt)
	if jwt != nil {
		jwtRaw := jwt.(string)
		if len(jwtRaw) > 0 {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwtRaw))
		}
	}

	traceId := ctx.Value(customclaims.ContextTraceId)
	if traceId != nil {
		traceIdRaw := traceId.(string)
		if len(traceIdRaw) > 0 {
			req.Header.Add("x-trace-id", traceIdRaw)
		}
	}

	if config.GetConfig().LogDebug {
		for name, values := range req.Header {
			// Loop over all values for the name.
			for _, value := range values {
				logger.Debug(ctx, fmt.Sprintf("%s: %v", name, value))
			}
		}
	}
}

func ExtractErrorMessageFromJsonBytes(ctx context.Context, data []byte, defaultMessage string) string {
	if len(data) == 0 {
		return defaultMessage
	}

	var e ErrorResponse
	body := string(data)
	logger.Debug(ctx, fmt.Sprintf("Error body:\n%s\n", body))
	err := json.Unmarshal(data, &e)
	if err != nil {
		logger.Error(ctx, err)
	} else {
		return e.Error
	}

	return defaultMessage

}
func HandleError(ctx context.Context, err error, w http.ResponseWriter) {
	statusCode := http.StatusInternalServerError
	w.Header().Add("Content-Type", "application/json;charset=UTF-8")
	if e, ok := err.(*customerrors.ServiceError); ok {
		statusCode = e.StatusCode
	}
	resp := ErrorResponse{Error: err.Error()}
	b, _ := MarshalFormat(ctx, resp)
	HandleResponse(w, b, statusCode)
}

func HandleResponse(w http.ResponseWriter, b []byte, statusCode int) {
	w.Header().Add("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(statusCode)
	w.Write(b)
}

func MarshalFormat(ctx context.Context, v interface{}) ([]byte, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logger.Error(ctx, err)
	} else {
		logger.Debug(ctx, string(b))
	}
	return b, err
}

func ProcessResponse(ctx context.Context, v interface{}, w http.ResponseWriter, statusCode int) {
	b, err := MarshalFormat(ctx, v)
	if err != nil {
		// Handle as a server error?
		HandleError(ctx, err, w)
		return
	}

	HandleResponse(w, b, statusCode)
}

// Validates that the response code is a 2xx
func IsOkResponse(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
