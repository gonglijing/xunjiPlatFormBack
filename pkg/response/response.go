package response

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

// JsonRes 数据返回通用JSON数据结构
type JsonRes struct {
	Code    int         `json:"code"`    // 错误码((0:成功, 1:失败, >1:错误码))
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 返回数据(业务接口定义具体数据结构)
	//Redirect string      `json:"redirect"` // 引导客户端跳转到指定路由
}

func firstResponseData(data []interface{}, fallback interface{}) interface{} {
	if len(data) > 0 {
		return data[0]
	}
	return fallback
}

func writeAttachment(r *ghttp.Request, fileName string, content io.ReadSeeker, contentType string) {
	r.Response.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	r.Response.Writer.Header().Add("Content-Type", contentType)
	http.ServeContent(r.Response.Writer, r.Request, fileName, time.Now(), content)
}

// Json 返回标准JSON数据。
func Json(r *ghttp.Request, code int, message string, data ...interface{}) {
	responseData := firstResponseData(data, g.Map{})
	r.Response.WriteJson(JsonRes{
		Code:    code,
		Message: message,
		Data:    responseData,
	})
}

// JsonExit 返回标准JSON数据并退出当前HTTP执行函数。
func JsonExit(r *ghttp.Request, code int, message string, data ...interface{}) {
	Json(r, code, message, data...)
	r.ExitAll()
}

// JsonRedirect 返回标准JSON数据引导客户端跳转。
func JsonRedirect(r *ghttp.Request, code int, message, redirect string, data ...interface{}) {
	responseData := firstResponseData(data, nil)
	r.Response.WriteJson(JsonRes{
		Code:    code,
		Message: message,
		Data:    responseData,
		//Redirect: redirect,
	})
}

// JsonRedirectExit 返回标准JSON数据引导客户端跳转，并退出当前HTTP执行函数。
func JsonRedirectExit(r *ghttp.Request, code int, message, redirect string, data ...interface{}) {
	JsonRedirect(r, code, message, redirect, data...)
	r.ExitAll()
}

// ToXls 向前端返回Excel文件 参数 content 为上面生成的io.ReadSeeker， fileTag 为返回前端的文件名
func ToXls(r *ghttp.Request, content io.ReadSeeker, fileTag string) {
	fileName := fmt.Sprintf("%s%s%s.xlsx", gtime.Now().String(), `-`, fileTag)
	writeAttachment(r, fileName, content, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	r.ExitAll()
}

// ToPlainText 输出流
func ToPlainText(r *ghttp.Request, content []byte, fileName string) {
	r.Response.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	r.Response.Writer.Header().Add("Content-Type", "text/plain;charset=UTF-8")
	r.Response.Write(content)
	r.Response.Flush()
}

// ToJsonFIle 向前端返回文件 参数 content 为上面生成的io.ReadSeeker， fileTag 为返回前端的文件名
func ToJsonFIle(r *ghttp.Request, content io.ReadSeeker, fileTag string) {
	ToJsonFile(r, content, fileTag)
}

// ToJsonFile 向前端返回文件 参数 content 为上面生成的io.ReadSeeker， fileTag 为返回前端的文件名
func ToJsonFile(r *ghttp.Request, content io.ReadSeeker, fileTag string) {
	fileName := fmt.Sprintf("%s%s%s.json", gtime.Now().Format("20060102150405"), `-`, fileTag)
	writeAttachment(r, fileName, content, "application/json")
	r.ExitAll()
}
