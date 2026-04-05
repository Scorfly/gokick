package gokick

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func kickErrorFromResponse(statusCode int, responseBody []byte) error {
	var apiErr errorResponse
	apiUnmarshalErr := json.Unmarshal(responseBody, &apiErr)
	if apiUnmarshalErr != nil {
		return fmt.Errorf(
			"failed to unmarshal error response (KICK status code: %d and body %q): %v",
			statusCode,
			string(responseBody),
			apiUnmarshalErr,
		)
	}
	if apiErr.Message != "" {
		return NewError(statusCode, apiErr.Message)
	}

	var oauthErr authErrorResponse
	oauthUnmarshalErr := json.Unmarshal(responseBody, &oauthErr)
	if oauthUnmarshalErr != nil {
		return fmt.Errorf(
			"failed to unmarshal error response (KICK status code: %d and body %q): %v",
			statusCode,
			string(responseBody),
			oauthUnmarshalErr,
		)
	}
	msg := oauthErr.Message
	if msg == "" {
		msg = oauthErr.Error
	}
	if msg != "" {
		e := NewError(statusCode, msg)
		if oauthErr.ErrorDescription != "" {
			e = e.WithDescription(oauthErr.ErrorDescription)
		}
		return e
	}

	return fmt.Errorf(
		"failed to unmarshal error response (KICK status code: %d and body %q): empty error message",
		statusCode,
		string(responseBody),
	)
}

func makeRequest[T any](
	ctx context.Context,
	request *Client,
	method string,
	path string,
	statusCode int,
	body io.Reader,
) (Response[T], error) {
	return makeRequestWithBaseURL[T](ctx, request, request.options.APIBaseURL, method, path, statusCode, body)
}

func makeRequestWithBaseURL[T any](
	ctx context.Context,
	request *Client,
	baseURL string,
	method string,
	path string,
	statusCode int,
	body io.Reader,
) (Response[T], error) {
	url := request.buildURL(baseURL, path)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return Response[T]{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := request.do(req)
	if err != nil {
		return Response[T]{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return Response[T]{}, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response[T]{}, fmt.Errorf("failed to read response body (KICK status code %d): %v", resp.StatusCode, err)
	}

	if resp.StatusCode != statusCode {
		return Response[T]{}, kickErrorFromResponse(resp.StatusCode, responseBody)
	}

	type successResponse struct {
		Result T `json:"data"`
	}

	var success successResponse

	err = json.Unmarshal(responseBody, &success)
	if err != nil {
		return Response[T]{}, fmt.Errorf(
			"failed to unmarshal response body (KICK status code %d and body %q): %v", resp.StatusCode, string(responseBody), err,
		)
	}

	return Response[T](success), nil
}

func makePaginatedRequest[T any](
	ctx context.Context,
	request *Client,
	method string,
	path string,
	statusCode int,
	body io.Reader,
) (PaginatedResponse[T], error) {
	return makePaginatedRequestWithBaseURL[T](ctx, request, request.options.APIBaseURL, method, path, statusCode, body)
}

func makePaginatedRequestWithBaseURL[T any](
	ctx context.Context,
	request *Client,
	baseURL string,
	method string,
	path string,
	statusCode int,
	body io.Reader,
) (PaginatedResponse[T], error) {
	url := request.buildURL(baseURL, path)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return PaginatedResponse[T]{}, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := request.do(req)
	if err != nil {
		return PaginatedResponse[T]{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return PaginatedResponse[T]{}, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return PaginatedResponse[T]{}, fmt.Errorf("failed to read response body (KICK status code %d): %v", resp.StatusCode, err)
	}

	if resp.StatusCode != statusCode {
		return PaginatedResponse[T]{}, kickErrorFromResponse(resp.StatusCode, responseBody)
	}

	var success PaginatedResponse[T]

	err = json.Unmarshal(responseBody, &success)
	if err != nil {
		return PaginatedResponse[T]{}, fmt.Errorf(
			"failed to unmarshal response body (KICK status code %d and body %q): %v", resp.StatusCode, string(responseBody), err,
		)
	}

	return success, nil
}

func makeAuthRequest[T any](
	ctx context.Context,
	request *Client,
	method string,
	path string,
	statusCode int,
	body io.Reader,
) (T, error) {
	url := request.buildURL(request.options.AuthBaseURL, path)

	var response T
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return response, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := request.do(req)
	if err != nil {
		return response, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return response, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read response body (KICK status code %d): %v", resp.StatusCode, err)
	}

	if resp.StatusCode != statusCode {
		var errorOutput authErrorResponse

		err = json.Unmarshal(responseBody, &errorOutput)
		if err != nil {
			return response, fmt.Errorf(
				"failed to unmarshal error response (KICK status code: %d and body %q): %v",
				resp.StatusCode,
				string(responseBody),
				err,
			)
		}

		if errorOutput.Message != "" {
			return response, NewError(resp.StatusCode, errorOutput.Message)
		} else {
			return response, NewError(resp.StatusCode, errorOutput.Error).WithDescription(errorOutput.ErrorDescription)
		}
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return response, fmt.Errorf(
			"failed to unmarshal response body (KICK status code %d and body %q): %v", resp.StatusCode, string(responseBody), err,
		)
	}

	return response, nil
}
