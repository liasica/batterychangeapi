package api

import (
    filetool "battery/library/file"
    "crypto/md5"
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"
    "regexp"
    "strings"
    "time"

    "battery/app/model"
    "battery/library/response"
)

var Upload = uploadApi{}

type uploadApi struct{}

// Image
// @Summary 公用-单图片上传
// @Tags    公用
// @Param   image formData file true "jpg/png图片"
// @Router  /api/upload/image [POST]
// @Success 200 {object} response.JsonResponse{data=model.UploadRep}  "返回结果"
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
        g.Log().Error(err.Error())
        response.JsonErrExit(r, response.RespCodeSystemError)
    }
    response.JsonOkExit(r, model.UploadRep{
        Path: fmt.Sprintf(fmt.Sprintf("%s%d%d%d/%s", "/uploads/", now.Year(), now.Month(), now.Day(), newName)),
    })
}

// Base64Image
// @Summary 公用-单图片上传(base64)
// @Tags    公用
// @Param   entity  body model.ImageBase64Req true "请求数据"
// @Router  /api/upload/base64_image [POST]
// @Success 200 {object} response.JsonResponse{data=model.UploadRep}  "返回结果"
func (*uploadApi) Base64Image(r *ghttp.Request) {
    var req model.ImageBase64Req
    if err := r.Parse(&req); err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }
    pattern := `^data:\s*image\/(\w+);base64,`

    b, _ := regexp.MatchString(pattern, req.Base64Content)
    if !b {
        response.Json(r, response.RespCodeArgs, "请上传base64文件")
    }
    re, _ := regexp.Compile(pattern)
    data := re.FindAllSubmatch([]byte(req.Base64Content), 2)
    fileType := string(data[0][1])
    fileData := re.ReplaceAllString(req.Base64Content, "")

    m := md5.New()
    m.Write([]byte(req.Base64Content))
    contentMd5 := hex.EncodeToString(m.Sum(nil))
    fileName := fmt.Sprintf("%s.%s", contentMd5, fileType)

    fileDataBin, err := base64.StdEncoding.DecodeString(fileData)
    if err != nil {
        response.JsonErrExit(r, response.RespCodeSystemError)
    }
    now := time.Now()
    dir := fmt.Sprintf("%s%d%d%d/%s", "./uploads/", now.Year(), now.Month(), now.Day(), contentMd5[30:])
    if _, err := os.Stat(dir); err != nil {
        _ = os.MkdirAll(dir, 0755)
    }
    err = ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, fileName), fileDataBin, 0666)
    if err != nil {
        response.JsonErrExit(r, response.RespCodeSystemError)
    }
    response.JsonOkExit(r, model.UploadRep{
        Path: fmt.Sprintf(fmt.Sprintf("%s%d%d%d/%s/%s", "/uploads/", now.Year(), now.Month(), now.Day(), contentMd5[30:], fileName)),
    })
}

// File
// @Summary 文件上传
// @Tags    公用
// @Param   file formData file true "文件"
// @Router  /api/upload/file [POST]
// @Success 200 {object} response.JsonResponse{data=model.UploadRep}  "返回结果"
func (*uploadApi) File(r *ghttp.Request) {
    file := r.GetUploadFile("file")
    if file == nil {
        response.Json(r, response.RespCodeArgs, "文件为空")
    }

    tool := &filetool.MultipartFile{FileHeader: file.FileHeader}
    e, err := tool.GetEtag()
    if err != nil {
        response.Json(r, response.RespCodeArgs, err.Error())
    }

    mimeType, err := tool.GetFileContentType()
    if err != nil {
        response.Json(r, response.RespCodeArgs, "未识别的文件")
    }

    ext := mimeType.Extension()
    dir := "./uploads/"
    if _, err = os.Stat(dir); err != nil {
        _ = os.MkdirAll(dir, 0755)
    }
    _path, err := file.Save(dir, true)
    if err != nil {
        g.Log().Error(err.Error())
        response.JsonErrExit(r, response.RespCodeSystemError)
    }

    dst := filepath.Join(dir, e+ext)
    _ = os.Rename(filepath.Join(dir, _path), dst)
    response.JsonOkExit(r, model.UploadRep{
        Path: dst,
    })
}
