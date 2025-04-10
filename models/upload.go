package models

// ChunkUpload 分片上传请求参数
type ChunkUpload struct {
	ChunkNumber int    `form:"chunkNumber" binding:"required"`
	TotalChunks int    `form:"totalChunks" binding:"required"`
	FileName    string `form:"fileName" binding:"required"`
	FileHash    string `form:"fileHash" binding:"required"`
}

// MergeRequest 合并文件请求参数
type MergeRequest struct {
	FileName    string `json:"fileName" binding:"required"`
	TotalChunks int    `json:"totalChunks" binding:"required"`
	FileHash    string `json:"fileHash" binding:"required"`
}