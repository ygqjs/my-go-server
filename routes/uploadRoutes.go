package routes

import (
	"my-go-server/controllers/upload"

	"github.com/gin-gonic/gin"
)

func UploadRoutes(r *gin.Engine) {
	uploadRoutes := r.Group("/upload")
	{
		// 上传分片
		uploadRoutes.POST("/chunk", upload.UploadController{}.UploadChunk)
		// 合并分片
		uploadRoutes.POST("/merge", upload.UploadController{}.MergeFile)
		// 查询已上传的分片状态
		uploadRoutes.GET("/status", upload.UploadController{}.GetUploadStatus)
	}
}