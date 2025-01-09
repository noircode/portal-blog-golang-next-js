package pagination

import (
	"math"
	"portal-blog/internal/core/domain/entity"
)

type IPagination interface {
	AddPagination(totalData, page, perPage int) (*entity.Page, error)
}

type Options struct{}

// AddPagination calculates pagination details based on the total number of items, requested page, and items per page.
//
// Parameters:
//   - totalData: The total number of items in the dataset.
//   - page: The requested page number (1-based index).
//   - perPage: The number of items to display per page.
//
// Returns:
//   - *entity.Page: A pointer to a Page struct containing pagination details.
//   - error: An error if the requested page is invalid or exceeds the maximum page count.
func (o *Options) AddPagination(totalData int, page int, perPage int) (*entity.Page, error) {
    newPage := page

    if newPage <= 0 {
        return nil, ErrorPage
    }

    limitData := 10

    if perPage > 0 {
        limitData = perPage
    }

    totalPage := int(math.Ceil(float64(totalData) / float64(limitData)))

    last := (newPage * limitData)
    first := last - limitData

    if totalData < last {
        last = totalData
    }

    zeroPage := &entity.Page{PageCount: 1, Page: newPage}

    if totalData == 0 && newPage == 1 {
        return zeroPage, nil
    }

    if newPage > totalPage {
        return nil, ErrorMaxPage
    }

    pages := &entity.Page{
        Page:     newPage,
    PerPage:  perPage,
    PageCount: totalPage,
    TotalCount: totalData,
    First:    first,
    Last:     last,
    }

    return pages, nil
}

// NewPagination creates and returns a new instance of IPagination.
//
// This function initializes a new Options struct and returns it as an IPagination interface.
// It serves as a factory method for creating pagination objects.
//
// Returns:
//   - IPagination: An interface that provides pagination functionality through the Options struct.
func NewPagination() IPagination {
    pagination := new(Options)

    return pagination
}
