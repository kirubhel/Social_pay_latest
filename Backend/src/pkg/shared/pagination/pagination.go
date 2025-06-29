package pagination

import (
	"math"

	"github.com/gin-gonic/gin"
	"github.com/socialpay/socialpay/src/pkg/shared/logging"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 10
	MaxPerPage     = 100
)

type Pagination struct {
	Page    int `json:"page" form:"page" binding:"required,min=1"`
	PerPage int `json:"page_size" form:"page_size" binding:"required,min=1,max=20"`
}


type PaginationInfo struct {
	Pagination 
	TotalItems int `json:"total_items"`
	TotalPage int `json:"total_page"`
	HasNextPage bool `json:"has_next_page"`
	HasPerviousPage bool `json:"has_pervious"`
}

func NewPagination(c *gin.Context, log logging.Logger) (*Pagination, error) {

	var p Pagination

	// Binding Pagination
	if err := c.ShouldBindQuery(&p); err != nil {

		log.Error("pagination error ", map[string]interface{}{
			"type":  "BINDING PAGINATION",
			"error": err,
		})
		return nil, err

	}

	if p.PerPage > MaxPerPage {
		p.PerPage = MaxPerPage
	}

	if p.Page < 1 {
		p.Page = DefaultPage
	}

	// returning pagination obj
	return &p, nil
}
// offet
func (p *Pagination) GetOffset() int {

	offset := (p.Page - 1) * p.PerPage

	return offset

}

//get limit
func (p *Pagination) GetLimit() int {
	return p.PerPage
}


// get pagination info 


func (p *Pagination)GetInfo(totalItems int) PaginationInfo {

	// total item is the len of data 

	// totalPages:=int(math.Ceil(float64(totalItems))/float64(p.PerPage))

	totalPages := int(math.Ceil(float64(totalItems) / float64(p.PerPage)))

	return PaginationInfo{
		Pagination: *p,
		TotalItems: totalItems,
		TotalPage: totalPages,
		HasNextPage: p.Page < totalPages,
		HasPerviousPage: p.Page > 1,
	}
}