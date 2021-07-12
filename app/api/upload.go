package api

import (
	"battery/app/model"
	"battery/library/response"
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"os"
	"path"
	"strings"
	"time"
)

var Upload = uploadApi{}

type uploadApi struct{}

// Image
// @summary 公用-单图片上传
// @tags    公用
// @param image formData file true "jpg/png图片"
// @router  /api/upload/image [POST]
// @success 200 {object} response.JsonResponse{data=model.UploadImageRep}  "返回结果"
func (*uploadApi) Image(r *ghttp.Request) {
	file := r.GetUploadFile("image")
	if file == nil {
		response.Json(r, response.RespCodeArgs, "文件为空")
	}
	ext := strings.ToLower(path.Ext(file.Filename))
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		response.Json(r, response.RespCodeArgs, "只支持jpg/png图片上传")
	}
	f, _ := file.Open()
	defer f.Close()
	now := time.Now()
	dir := fmt.Sprintf("%s%d%d%d", "./uploads/", now.Year(), now.Month(), now.Day())
	if _, err := os.Stat(dir); err != nil {
		_ = os.MkdirAll(dir, 0755)
	}
	newName, err := file.Save(dir, true)
	if err != nil {
		response.JsonErrExit(r, response.RespCodeSystemError)
	}
	response.JsonOkExit(r, model.UploadImageRep{
		Path: fmt.Sprintf(fmt.Sprintf("%s%d%d%d/%s", "/uploads/", now.Year(), now.Month(), now.Day(), newName)),
	})
}
