package gokick

import (
	"context"
	"fmt"
	"net/http"
)

type (
	CategoriesResponseWrapper Response[[]CategoryResponse]
	CategoryResponseWrapper   Response[CategoryResponse]
)

type CategoryResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}

func (c *Client) GetCategories(ctx context.Context) (CategoriesResponseWrapper, error) {
	response, err := makeRequest[[]CategoryResponse](ctx, c, http.MethodGet, "/public/v1/categories", http.StatusOK, http.NoBody)
	if err != nil {
		return CategoriesResponseWrapper{}, err
	}

	return CategoriesResponseWrapper(response), nil
}

func (c *Client) GetCategory(ctx context.Context, categoryID int) (CategoryResponseWrapper, error) {
	response, err := makeRequest[[]CategoryResponse](
		ctx,
		c,
		http.MethodGet,
		fmt.Sprintf("/public/v1/categories/%d", categoryID),
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return CategoryResponseWrapper{}, err
	}

	return CategoryResponseWrapper{
		Result: response.Result[0],
	}, nil
}
