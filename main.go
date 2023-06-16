package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	router := gin.Default()

	router.POST("/converter/m3u8", func(c *gin.Context) {
		//// 获取当前项目的绝对路径
		projectPath, err := filepath.Abs(".")
		if err != nil {
			log.Println("获取项目路径失败:", err)
			c.String(http.StatusInternalServerError, "转码失败")
			return
		}
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
		}
		// 获取所有图片
		files := form.File["files"]
		log.Println("files=>", files)
		// 遍历所有图片
		for _, file := range files {
			// 逐个存
			fileName := strings.Split(file.Filename, ".")[0]
			savePath := filepath.Join(projectPath, "static", fileName, file.Filename)
			if err := c.SaveUploadedFile(file, savePath); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
				return
			}

			//构建输入文件和输出文件的绝对路径
			outputPath := filepath.Join(projectPath, "static", fileName, fileName+".m3u8")

			// 调用 FFmpeg 命令行进行转码
			cmd := exec.Command("pkg/ffmpeg", "-i", savePath, outputPath)
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Println("转码失败:", err)
				log.Println("FFmpeg 输出:", string(output))
				c.String(http.StatusInternalServerError, "转码失败")
				return
			}
		}

		c.String(http.StatusOK, "转码成功")
	})

	router.Run(":8080")
}
