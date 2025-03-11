package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"logicflow-deploy/internal/schema"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	fileNameRegex = regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	flowDataPath  = "data/flowdata"
)

func CreateFlowData(c *gin.Context) {
	var data schema.FlowData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request data"})
		return
	}

	// 从url中获取文件名
	name := c.Param("name")

	// 验证文件名是否符合要求
	if !fileNameRegex.MatchString(name) {
		c.JSON(400, gin.H{"error": "Invalid file name"})
		return
	}

	// 检查文件是否已存在
	filePath := filepath.Join(flowDataPath, name+".json")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		c.JSON(409, gin.H{"error": "File already exists"})
		return
	}
	// 保存文件
	if err := saveFlowDataToFile(data, filePath); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}
	c.JSON(201, gin.H{"message": "文件创建成功", "filename": name})
}

func saveFlowDataToFile(data schema.FlowData, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// 获取Flow文件详情
func GetFlow(c *gin.Context) {
	filename := c.Param("name")
	if !validateFilename(filename) {
		c.JSON(400, gin.H{"error": "无效的文件名"})
		return
	}

	content, err := readFlowFile(filepath.Join(flowDataPath, filename+".json"))
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(404, gin.H{"error": "文件不存在"})
		} else {
			c.JSON(500, gin.H{"error": "文件读取失败"})
		}
		return
	}

	c.JSON(200, content)
}

// 更新Flow文件
func UpdateFlow(c *gin.Context) {
	filename := c.Param("filename")
	if !validateFilename(filename) {
		c.JSON(400, gin.H{"error": "无效的文件名"})
		return
	}

	var content map[string]interface{}
	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求内容"})
		return
	}

	filePath := filepath.Join(flowDataPath, filename+".json")
	if err := writeFlowFile(filePath, content); err != nil {
		c.JSON(500, gin.H{"error": "文件更新失败"})
		return
	}

	c.JSON(200, gin.H{"message": "文件更新成功"})
}

// 删除Flow文件（通过重命名标记删除）
func DeleteFlow(c *gin.Context) {
	filename := c.Param("filename")
	if !validateFilename(filename) {
		c.JSON(400, gin.H{"error": "无效的文件名"})
		return
	}

	oldPath := filepath.Join(flowDataPath, filename+".json")
	newPath := filepath.Join(flowDataPath, filename+".deleted")

	if err := os.Rename(oldPath, newPath); err != nil {
		if os.IsNotExist(err) {
			c.JSON(404, gin.H{"error": "文件不存在"})
		} else {
			c.JSON(500, gin.H{"error": "文件删除失败"})
		}
		return
	}

	c.JSON(200, gin.H{"message": "文件已标记删除"})
}

// 辅助函数：验证文件名格式
func validateFilename(name string) bool {
	return fileNameRegex.MatchString(name)
}

// 辅助函数：读取文件内容
func readFlowFile(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content map[string]interface{}
	err = json.NewDecoder(file).Decode(&content)
	return content, err
}

// 辅助函数：写入文件
func writeFlowFile(path string, content interface{}) error {
	// 确保目录存在
	if err := os.MkdirAll(flowDataPath, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(content)
}

// 返回所有Flow文件名列表
func ListFlow(c *gin.Context) {
	// 列出所有Flow文件
	// 这里需要根据实际情况实现
	// 假设flowDataPath是Flow文件所在的目录
	files, err := os.ReadDir(flowDataPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "读取文件列表失败"})
		return
	}
	type FileInfo struct {
		Name      string    `json:"name"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
	allFiles := make([]FileInfo, 0)
	for _, f := range files {
		if !f.IsDir() {
			fileInfo, _ := f.Info()
			name := fileInfo.Name()
			allFiles = append(allFiles, FileInfo{
				Name:      strings.Split(name, ".json")[0],
				UpdatedAt: fileInfo.ModTime(),
			})
		}
	}
	c.JSON(200, gin.H{"message": "列出所有Flow文件",
		"data": allFiles})
}
