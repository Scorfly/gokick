package gokick

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func makeRequest[T any](
	ctx context.Context,
	request *Client,
	method string,
	path string,
	statusCode int,
	body io.Reader,
) (Response[T], error) {
	url := request.buildURL(request.options.APIBaseURL, path)

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
		var errorOutput errorResponse

		err = json.Unmarshal(responseBody, &errorOutput)
		if err != nil {
			return Response[T]{}, fmt.Errorf(
				"failed to unmarshal error response (KICK status code: %d and body %q): %v",
				resp.StatusCode,
				string(responseBody),
				err,
			)
		}

		return Response[T]{}, NewError(resp.StatusCode, errorOutput.Message)
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
