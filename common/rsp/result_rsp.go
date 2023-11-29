package rsp

import (
	"errors"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"

	//"encoding/json"
	"net/http"
)

type ResultRsp struct {
	code int         `json:"code"`
	msg  string      `json:"success"`
	data interface{} `json:data"`
}

// 正确状态处理
func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  msg,
		"data": data,
	})
}

// 错误状态处理
func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest,
		gin.H{
			"error": msg,
		})
}

//============================================================

var (
	// locale2FieldMap map[string]map[string]string

	// locale = "en"

	translator ut.Translator
)

// ShouldBindJSON 取得 Request body 資料，JSON 資料轉行為 struct
func ShouldBindJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		var messages string
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, err2 := range err.(validator.ValidationErrors) {
				messages += err2.Translate(translator) + ", "
			}

			// 過濾最後字尾 ", "
			messages = messages[:len(messages)-2]

		} else {
			messages = err.Error()
		}

		return errors.New(messages)
	}

	return nil
}

type Basic struct {
	Status Status `json:"status"`

	Data interface{} `json:"data"`
}

func NewSuccess(data interface{}) Basic {
	return Basic{
		Status: CodeOk,
		Data:   data,
	}
}

func NewError(status Status) Basic {
	return Basic{
		Status: status,
		Data:   []string{},
	}
}
