package initialize

import (
	"mxshop_api/user-web/global"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
)

func InitTranslations() {
	// 获取gin的验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册中文翻译
		zhs := zh.New()
		uni := ut.New(zhs, zhs)
		trans, _ := uni.GetTranslator("zh")

		// 注册默认中文翻译
		_ = zh_trans.RegisterDefaultTranslations(v, trans)

		// 赋值给全局变量
		global.Trans = trans
	}
}
