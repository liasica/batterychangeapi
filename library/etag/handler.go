// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-14
// Based on apiv2 by liasica, magicrolan@qq.com.

package etag

import (
    "bytes"
    "crypto/sha1"
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "github.com/gogf/gf/net/ghttp"
    "hash"
    "log"
    "net/http"
    "strconv"

    "github.com/go-http-utils/fresh"
    "github.com/go-http-utils/headers"
)

type hashWriter struct {
    rw     http.ResponseWriter
    hash   hash.Hash
    buf    *bytes.Buffer
    len    int
    status int
}

func (hw hashWriter) Header() http.Header {
    return hw.rw.Header()
}

func (hw *hashWriter) WriteHeader(status int) {
    hw.status = status
}

func (hw *hashWriter) Write(b []byte) (int, error) {
    if hw.status == 0 {
        hw.status = http.StatusOK
    }
    // bytes.Buffer.Write(b) always return (len(b), nil), so just
    // ignore the return values.
    hw.buf.Write(b)

    l, err := hw.hash.Write(b)
    hw.len += l
    return l, err
}

func writeRaw(res http.ResponseWriter, hw hashWriter) {
    res.WriteHeader(hw.status)
    res.Write(hw.buf.Bytes())
}

func Handler(r *ghttp.Request, weak bool) {
    h := sha1.New()
    h.Write(r.Response.Buffer())
    etag := fmt.Sprintf("%X-%v", r.Response.BufferLength(), base64.StdEncoding.EncodeToString(h.Sum(nil)))
    if weak {
        etag = "W/" + etag
    }

    log.Println(r.Request.Header.Get(headers.ETag))

    if fresh.IsFresh(r.Request.Header, r.Response.Header()) {
        log.Println("xxx")
        r.Response.WriteHeader(http.StatusNotModified)
        r.Response.Write(nil)
    } else {
        r.Response.Header().Set(headers.ETag, etag)
    }
}

// HttpHandler wraps the http.Handler h with ETag support.
func HttpHandler(h http.Handler, weak bool) http.Handler {
    return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        hw := hashWriter{rw: res, hash: sha1.New(), buf: bytes.NewBuffer(nil)}
        h.ServeHTTP(&hw, req)

        resHeader := res.Header()

        if hw.hash == nil ||
            resHeader.Get(headers.ETag) != "" ||
            strconv.Itoa(hw.status)[0] != '2' ||
            hw.status == http.StatusNoContent ||
            hw.buf.Len() == 0 {
            writeRaw(res, hw)
            return
        }

        etag := fmt.Sprintf("%v-%v", strconv.Itoa(hw.len),
            hex.EncodeToString(hw.hash.Sum(nil)))

        if weak {
            etag = "W/" + etag
        }

        resHeader.Set(headers.ETag, etag)

        if fresh.IsFresh(req.Header, resHeader) {
            res.WriteHeader(http.StatusNotModified)
            res.Write(nil)
        } else {
            writeRaw(res, hw)
        }
    })
}
