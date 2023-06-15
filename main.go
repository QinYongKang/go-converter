package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
)

func main() {
	router := gin.Default()

	router.POST("/mp4-to-m3u8", func(c *gin.Context) {
		//file, err := c.FormFile("mp4")
		//if err != nil {
		//	log.Println("err==>", err)
		//	c.String(500, "上传文件出错")
		//	return
		//}
		//// c.JSON(200, gin.H{"message": file.Header.Context})
		//log.Println("上传完毕，准备保存==>")
		//
		//// 获取当前项目的绝对路径
		projectPath, err := filepath.Abs(".")
		if err != nil {
			log.Println("获取项目路径失败:", err)
			c.String(http.StatusInternalServerError, "转码失败")
			return
		}
		//savePath := filepath.Join(projectPath, "static", file.Filename)
		//c.SaveUploadedFile(file, savePath)
		//c.String(http.StatusOK, file.Filename)

		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get err %s", err.Error()))
		}
		// 获取所有图片
		files := form.File["files"]
		// 遍历所有图片
		for _, file := range files {
			// 逐个存
			savePath := filepath.Join(projectPath, "static", file.Filename)
			if err := c.SaveUploadedFile(file, savePath); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload err %s", err.Error()))
				return
			}

			//构建输入文件和输出文件的绝对路径
			//inputPath := filepath.Join(projectPath, "static", "input.mp4")
			//outputPath := filepath.Join(projectPath, "static", "output.m3u8")

			// 调用 FFmpeg 命令行进行转码
			cmd := exec.Command("pkg/ffmpeg", "-i", savePath, file.Filename+".m3u8")
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Println("转码失败:", err)
				log.Println("FFmpeg 输出:", string(output))
				c.String(http.StatusInternalServerError, "转码失败")
				return
			}
		}

		//err := cmd.Run()
		//if err != nil {
		//	log.Println("转码失败:", err)
		//	c.String(http.StatusInternalServerError, "转码失败")
		//	return
		//}

		c.String(http.StatusOK, "转码成功")
	})

	router.Run(":8080")
}
