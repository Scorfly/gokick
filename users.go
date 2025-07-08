package gokick

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type (
	UsersResponseWrapper           Response[[]UserResponse]
	UserResponseWrapper            Response[UserResponse]
	TokenIntrospectResponseWrapper Response[TokenIntrospectResponse]
)

type UserResponse struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
	UserID         int    `json:"user_id"`
}

type TokenIntrospectResponse struct {
	Active    bool   `json:"active"`
	ClientID  string `json:"client_id"`
	Exp       int    `json:"exp"`
	Scope     string `json:"scope"`
	TokenType string `json:"token_type"`
}

type UserListFilter struct {
	queryParams url.Values
}

func NewUserListFilter() UserListFilter {
	return UserListFilter{queryParams: make(url.Values)}
}

func (f UserListFilter) SetID(id int) UserListFilter {
	f.queryParams.Set("id", fmt.Sprintf("%d", id))
	return f
}

func (f UserListFilter) SetIDs(ids []int) UserListFilter {
	for i := range ids {
		f.queryParams.Add("id", fmt.Sprintf("%d", ids[i]))
	}

	return f
}

func (f UserListFilter) ToQueryString() string {
	if len(f.queryParams) == 0 {
		return ""
	}

	return "?" + f.queryParams.Encode()
}

func (c *Client) TokenIntrospect(ctx context.Context) (TokenIntrospectResponseWrapper, error) {
	response, err := makeRequest[TokenIntrospectResponse](
		ctx,
		c,
		http.MethodPost,
		"/public/v1/token/introspect",
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return TokenIntrospectResponseWrapper{}, err
	}

	return TokenIntrospectResponseWrapper(response), nil
}

func (c *Client) GetUsers(ctx context.Context, filter UserListFilter) (UsersResponseWrapper, error) {
	response, err := makeRequest[[]UserResponse](
		ctx,
		c,
		http.MethodGet,
		fmt.Sprintf("/public/v1/users%s", filter.ToQueryString()),
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return UsersResponseWrapper{}, err
	}

	return UsersResponseWrapper(response), nil
}
