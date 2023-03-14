package request

import "rederinghub.io/internal/entity"

type ListDaoProjectRequest struct {
	entity.Pagination
	Status  *int64  `query:"status"`
	Keyword *string `query:"keyword"`
}
type CreateDaoProjectRequest struct {
	ProjectId string `json:"project_id"`
	CreatedBy string `json:"-"`
}