package gokick

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type (
	CategoriesResponseWrapper PaginatedResponse[[]CategoryResponse]
	CategoryResponseWrapper   Response[CategoryResponse]
)

type CategoryResponse struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Tags      []string `json:"tags,omitempty"`
	Thumbnail string   `json:"thumbnail"`
}

type CategoryListFilter struct {
	cursor string
	limit  *int
	names  []string
	tags   []string
	ids    []int
}

func NewCategoryListFilter() CategoryListFilter {
	return CategoryListFilter{}
}

// SetCursor sets the pagination cursor (v2).
func (f CategoryListFilter) SetCursor(cursor string) CategoryListFilter {
	f.cursor = cursor
	return f
}

// SetLimit sets the page size (1–1000, API default 25).
func (f CategoryListFilter) SetLimit(limit int) CategoryListFilter {
	f.limit = &limit
	return f
}

// AddName adds a category name filter (comma-serialized per Kick v2 form style).
func (f CategoryListFilter) AddName(name string) CategoryListFilter {
	f.names = append(f.names, name)
	return f
}

// AddTag adds a category tag filter.
func (f CategoryListFilter) AddTag(tag string) CategoryListFilter {
	f.tags = append(f.tags, tag)
	return f
}

// AddID adds a category id filter.
func (f CategoryListFilter) AddID(id int) CategoryListFilter {
	f.ids = append(f.ids, id)
	return f
}

func (f CategoryListFilter) ToQueryString() string {
	v := url.Values{}
	if f.cursor != "" {
		v.Set("cursor", f.cursor)
	}
	if f.limit != nil {
		v.Set("limit", strconv.Itoa(*f.limit))
	}
	if len(f.names) > 0 {
		v.Set("name", strings.Join(f.names, ","))
	}
	if len(f.tags) > 0 {
		v.Set("tag", strings.Join(f.tags, ","))
	}
	if len(f.ids) > 0 {
		idStrs := make([]string, len(f.ids))
		for i, id := range f.ids {
			idStrs[i] = strconv.Itoa(id)
		}
		v.Set("id", strings.Join(idStrs, ","))
	}
	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}

func (c *Client) GetCategories(ctx context.Context, filter CategoryListFilter) (CategoriesResponseWrapper, error) {
	response, err := makePaginatedRequest[[]CategoryResponse](
		ctx,
		c,
		http.MethodGet,
		fmt.Sprintf("/public/v2/categories%s", filter.ToQueryString()),
		http.StatusOK,
		http.NoBody,
	)
	if err != nil {
		return CategoriesResponseWrapper{}, err
	}

	return CategoriesResponseWrapper(response), nil
}

func (c *Client) GetCategory(ctx context.Context, categoryID int) (CategoryResponseWrapper, error) {
	filter := NewCategoryListFilter().AddID(categoryID).SetLimit(1)
	list, err := c.GetCategories(ctx, filter)
	if err != nil {
		return CategoryResponseWrapper{}, err
	}
	if len(list.Result) == 0 {
		return CategoryResponseWrapper{}, fmt.Errorf("category id %d: empty result", categoryID)
	}
	return CategoryResponseWrapper{Result: list.Result[0]}, nil
}
