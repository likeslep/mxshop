package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"mxshop_api/user-web/global" 
)

// TranslateError 翻译单条错误信息
func TranslateError(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			// 使用全局翻译器
			return e.Translate(global.Trans)
		}
	}
	return "参数校验失败"
}

// HandleValidatorError 表单验证统一返回函数
// 所有 handler 表单错误直接调用这一行即可
func HandleValidatorError(c *gin.Context, err error) {
	msg := TranslateError(err)
	c.JSON(http.StatusBadRequest, gin.H{
		"code": 400,
		"msg":  msg,
	})
}