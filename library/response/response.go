package response

import (
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
    "reflect"
)

const (
    RespCodeSuccess      = 200
    RespCodeArgs         = 400
    RespCodeUnauthorized = 401
    RespCodeNotFound     = 404
    RespCodeSystemError  = 500
)

var CodeMsg = map[int]string{
    RespCodeSuccess:      "成功",
    RespCodeArgs:         "参数错误",
    RespCodeUnauthorized: "请先登录",
    RespCodeNotFound:     "不存在的资源",
    RespCodeSystemError:  "系统繁忙，请稍后再试！",
}

// JsonResponse 数据返回通用JSON数据结构
type JsonResponse struct {
    Code    int         `validate:"required" json:"code"`    // 错误码((200:成功, 其它:失败))
    Message string      `validate:"required" json:"message"` // 提示信息
    Data    interface{} `json:"data,omitempty"`              // 返回数据(业务接口定义具体数据结构)
}

// Json 标准返回结果数据结构封装。
func Json(r *ghttp.Request, code int, message string, data ...interface{}) {
    responseData := interface{}(nil)
    if len(data) > 0 {
        responseData = data[0]
    }
    _ = r.Response.WriteJson(JsonResponse{
        Code:    code,
        Message: message,
        Data:    responseData,
    })
    r.Exit()
}

// JsonOkExit 返回json请求并退出
func JsonOkExit(r *ghttp.Request, data ...interface{}) {
    responseData := interface{}(nil)
    if len(data) > 0 {
        responseData = data[0]
    }
    _ = r.Response.WriteJson(JsonResponse{
        Code:    RespCodeSuccess,
        Message: CodeMsg[RespCodeSuccess],
        Data:    responseData,
    })
    r.Exit()
}

// JsonErrExit 返回json错误并退出
func JsonErrExit(r *ghttp.Request, args ...int) {
    rep := JsonResponse{
        Code:    RespCodeSystemError,
        Message: CodeMsg[RespCodeSystemError],
    }
    l := len(args)
    if l > 0 {
        rep.Code = args[0]
        if msg, ok := CodeMsg[args[0]]; ok {
            rep.Message = msg
        }
    }
    _ = r.Response.WriteJson(rep)
    r.Exit()
}

func ItemsWithTotal(r *ghttp.Request, total int, items interface{}) {
    v := reflect.ValueOf(items)
    if v.IsNil() || v.IsZero() {
        items = []interface{}{}
    }
    JsonOkExit(r, g.Map{
        "total": total,
        "items": items,
    })
}
