package Query

// PageQuery
// @Description: 分页入参
type PageQuery struct {
	PageNum  *int `json:"pageNum" form:"pageNum" binding:"required,min=1"`
	PageSize *int `json:"pageSize" form:"pageSize" binding:"required,min=1"`
}
