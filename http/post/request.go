package post

import "github.com/Artexus/api-matthew-backend/utils/pagination"

type InsertPostRequest struct {
	UserID      int     `json:"-"`
	FileID      *string `json:"file_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
}

type GetPostRequest struct {
	pagination.Pagination
}
