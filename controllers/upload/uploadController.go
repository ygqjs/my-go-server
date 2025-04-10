package upload

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"my-go-server/database"
	"my-go-server/models"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
}

/**
 * 查询已上传的分片状态
 * 接口：GET /upload/status?fileHash={fileHash}
 * 返回值示例：
 * {
 *   "success": true,
 *   "data": {
 *     "uploadedChunks": [1, 2, 3]
 *   },
 *   "message": "查询成功",
 * }
 */
func (uploadController UploadController) GetUploadStatus(c *gin.Context) {
	fileHash := c.Query("fileHash")
	if fileHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "fileHash is required",
		})
		return
	}

	// 使用 DB.Raw 执行原生 SQL 查询
	uploadedChunks := []int{}
	err := database.DB.Raw("SELECT chunk_number FROM file_chunks WHERE file_hash = ?", fileHash).Scan(&uploadedChunks).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to query uploaded chunks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"uploadedChunks": uploadedChunks,
		},
	})
}

/**
 * 处理分片上传
 * 接口：POST /upload/chunk
 * 请求参数（FormData）：
 * - chunk: 文件分片内容
 * - chunkNumber: 当前分片号
 * - totalChunks: 总分片数
 * - fileName: 文件名
 * - fileID: 文件哈希（前端传来的 fileHash）
 * 功能：将分片数据插入数据库
 * 返回值示例：
 * {
 *   "success": true,
 *   "data": {
 *     "chunkNumber": 1
 *   }
 *   "message": "分片上传成功",
 * }
 */
func (uploadController UploadController) UploadChunk(c *gin.Context) {
	// 获取分片文件
	file, _, err := c.Request.FormFile("chunk")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无法获取文件分片",
		})
		return
	}
	defer file.Close()

	// 解析表单参数
	var chunkUpload models.ChunkUpload
	if err := c.ShouldBind(&chunkUpload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数不完整",
		})
		return
	}

	fileHash := chunkUpload.FileHash

	// 读取分片内容
	chunkData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "读取分片失败",
		})
		return
	}

	// 使用原生 SQL 插入分片数据
	query := "INSERT INTO file_chunks (file_hash, file_name, chunk_number, total_chunks, chunk_data) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE chunk_data = VALUES(chunk_data)"
	result := database.DB.Exec(query, fileHash, chunkUpload.FileName, chunkUpload.ChunkNumber, chunkUpload.TotalChunks, chunkData)
	if result.Error != nil {
		log.Printf("Failed to insert chunk: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("保存分片到数据库失败: %v", result.Error),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"chunkNumber": chunkUpload.ChunkNumber,
		},
		"message": "分片上传成功",
	})
}

/**
 * 合并分片文件
 * 接口：POST /upload/merge
 * 请求参数：
 * - fileName: 文件名
 * - totalChunks: 总分片数
 * - fileHash: 文件哈希（前端传递的参数名已改为 fileHash）
 * 功能：从数据库读取分片，合并为文件，保存到本地并清理数据库记录
 * 返回值示例：
 * {
 *   "success": true,
 *   "data": {
 *     "fileName": "example.zip",
 *     "fileHash": "abc123..."
 *   }
 *   "message": "文件上传完成",
 * }
 */
func (uploadController UploadController) MergeFile(c *gin.Context) {
	var req models.MergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "参数不完整",
		})
		return
	}

	fileHash := req.FileHash

	// 检查分片数量是否齐全
	var count int
	err := database.DB.Raw("SELECT COUNT(*) FROM file_chunks WHERE file_hash = ? AND total_chunks = ?", fileHash, req.TotalChunks).Scan(&count).Error
	if err != nil || count != req.TotalChunks {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "分片不完整",
		})
		return
	}

	// 创建最终文件存储目录
	finalDir := "./final_files/"
	if err := os.MkdirAll(finalDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建目录失败",
		})
		return
	}

	// 分离文件名和扩展名，生成保存文件名
	ext := filepath.Ext(req.FileName)                                // 提取扩展名，例如 ".md"
	baseName := req.FileName[:len(req.FileName)-len(ext)]            // 提取基础文件名，例如 "codereview"
	finalFileName := fmt.Sprintf("%s_%s%s", baseName, fileHash, ext) // 保存为 "codereview_5d08ced39910341325c102af785beb54.md"
	finalFilePath := filepath.Join(finalDir, finalFileName)

	// 创建目标文件
	finalFile, err := os.Create(finalFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建目标文件失败",
		})
		return
	}
	defer finalFile.Close()

	// 查询并合并分片数据
	var chunks []struct {
		ChunkData []byte `gorm:"column:chunk_data"`
	}
	err = database.DB.Raw("SELECT chunk_data FROM file_chunks WHERE file_hash = ? ORDER BY chunk_number ASC", fileHash).Scan(&chunks).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询分片数据失败",
		})
		return
	}

	for _, chunk := range chunks {
		_, err = finalFile.Write(chunk.ChunkData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "合并文件失败",
			})
			return
		}
	}

	// 计算最终文件的 MD5
	hash := md5.New()
	finalFile.Seek(0, 0)
	_, err = io.Copy(hash, finalFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "计算MD5失败",
		})
		return
	}
	md5Sum := hex.EncodeToString(hash.Sum(nil))

	// 删除分片记录以释放数据库空间
	err = database.DB.Exec("DELETE FROM file_chunks WHERE file_hash = ?", fileHash).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "清理分片记录失败",
		})
		return
	}

	// 返回指定结构
	c.JSON(http.StatusOK, gin.H{
		"success": "success",
		"data": gin.H{
			"fileName": req.FileName, // 返回原始文件名 "codereview.md"
			"fileHash": fileHash,
			"md5":      md5Sum,
		},
		"message": "文件上传完成",
	})
}
