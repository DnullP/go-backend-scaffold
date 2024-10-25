package main

import (
	"go-backend-scaffold/config"
	"go-backend-scaffold/proto"
	"go-backend-scaffold/services/client"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的 Gin 引擎
	config.LoadConfig()
	router := gin.Default()

	// 定义一个简单的 GET 路由
	router.GET("/get_user", func(c *gin.Context) {
		req := proto.GetUserRequest{
			Id: 123,
		}

		res, _ := client.GetUser(&req)
		c.JSON(http.StatusOK, gin.H{
			"name":  res.Name,
			"email": res.Email,
		})
	})

	// 定义一个 POST 路由
	router.POST("/submit", func(c *gin.Context) {
		var json struct {
			Title   string `json:"title" binding:"required"`
			Content string `json:"content" binding:"required"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "提交成功",
			"title":   json.Title,
			"content": json.Content,
		})
	})

	// 启动服务器，监听在 8080 端口
	router.Run(":8080")
}
