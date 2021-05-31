package httptools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

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
	return doApiCall(ctx, "POST", url, requestModel)
}

// Does a PUT to the specified endpoint. Returns the body bytes, an http status code (0 if no call was made)
func Put(ctx context.Context, url string, requestModel interface{}) ([]byte, int, error) {
	return doApiCall(ctx, "PUT", url, requestModel)
}

// Does a PUT to the specified endpoint. Returns the body bytes, an http status code (0 if no call was made)
func Patch(ctx context.Context, url string, requestModel interface{}) ([]byte, int, error) {
	return doApiCall(ctx, "PATCH", url, requestModel)
}

func Get(ctx context.Context, url string) ([]byte, int, error) {
	return doApiCall(ctx, "GET", url, nil)
}

func Delete(ctx context.Context, url string) ([]byte, int, error) {
	return doApiCall(ctx, "DELETE", url, nil)
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

func doApiCall(ctx context.Context, method string, url string, requestModel interface{}) ([]byte, int, error) {
	// Convert the request model into JSON

	logger.Debug(ctx, fmt.Sprintf("Preparing call to %s %s", method, url))
	requestBytes := []byte{}
	if !strings.EqualFold(method, "DELETE") && !strings.EqualFold(method, "GET") {
		var err error
		requestBytes, err = MarshalFormat(ctx, requestModel)
		if err != nil {
			logger.Error(ctx, err)
			return nil, 0, err
		}
		logger.Debug(ctx, fmt.Sprintf("JSON for %s %s: %s", method, url, string(requestBytes)))
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		logger.Error(ctx, err)
		return nil, 0, err
	}

	// Set the Authorization and x-trace-id headers
	setRequestHeader(ctx, req)

	// Make the call
	client := &http.Client{}
	logger.Debug(ctx, fmt.Sprintf("Calling %s %s", method, url))
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(ctx, err)
		return nil, 0, err
	}
	logger.Debug(ctx, fmt.Sprintf("Response (%s) from %s %s", resp.Status, method, url))
	defer resp.Body.Close()

	// Grab the body

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(ctx, err)
		return nil, resp.StatusCode, err
	}
	if len(bodyBytes) > 0 {
		logger.Debug(ctx, fmt.Sprintf("Response (%s): %s", resp.Status, string(bodyBytes)))
	}

	if !IsOkResponse(resp) {
		// TODO: try to extract out an error from the body
		errorMessage := extractErrorMessageFromJsonBytes(ctx, bodyBytes, fmt.Sprintf("HTTP Exception calling %s %s returned status %s", method, url, resp.Status))

		return nil, resp.StatusCode, customerrors.NewServiceError(resp.StatusCode, errorMessage)
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
}

func extractErrorMessageFromJsonBytes(ctx context.Context, data []byte, defaultMessage string) string {
	if len(data) == 0 {
		return defaultMessage
	}

	var e ErrorResponse
	err := json.Unmarshal(data, &e)
	if err != nil {
		logger.Error(ctx, err)
	} else {
		return e.Error
	}

	return defaultMessage
}
