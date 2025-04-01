package gokick

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

type CategoryResponseV2 struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CategoryListFilter struct {
	queryParams url.Values
}

func NewCategoryListFilter() CategoryListFilter {
	return CategoryListFilter{queryParams: make(url.Values)}
}

func (f CategoryListFilter) SetQuery(query string) CategoryListFilter {
	f.queryParams.Set("q", query)
	return f
}

func (f CategoryListFilter) ToQueryString() string {
	if len(f.queryParams) == 0 {
		return ""
	}

	return "?" + f.queryParams.Encode()
}

func (c *Client) GetCategories(ctx context.Context, filter CategoryListFilter) (CategoriesResponseWrapper, error) {
	response, err := makeRequest[[]CategoryResponse](
		ctx,
		c,
		http.MethodGet,
		fmt.Sprintf("/public/v1/categories%s", filter.ToQueryString()),
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return CategoriesResponseWrapper{}, err
	}

	return CategoriesResponseWrapper(response), nil
}

func (c *Client) GetCategory(ctx context.Context, categoryID int) (CategoryResponseWrapper, error) {
	response, err := makeRequest[CategoryResponse](
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

	return CategoryResponseWrapper(response), nil
}
