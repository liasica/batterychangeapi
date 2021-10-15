// Copyright (C) liasica. 2021-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
// Created at 2021-10-15
// Based on apiv2 by liasica, magicrolan@qq.com.

package file

import (
    "bytes"
    "crypto/sha1"
    "encoding/base64"
    "io"
)

const (
    BlockBits = 22 // Indicate that the blocksize is 4M
    BlockSize = 1 << BlockBits
)

func BlockCount(fsize int64) int {

    return int((fsize + (BlockSize - 1)) >> BlockBits)
}

func CalSha1(b []byte, r io.Reader) ([]byte, error) {

    h := sha1.New()
    _, err := io.Copy(h, r)
    if err != nil {
        return nil, err
    }
    return h.Sum(b), nil
}

func (f *MultipartFile) GetEtag() (etag string, err error) {

    blockCnt := BlockCount(f.Size)
    sha1Buf := make([]byte, 0, 21)

    r, err := f.Open()
    if err != nil {
        return "", err
    }
    defer r.Close()

    if blockCnt <= 1 { // file size <= 4M
        sha1Buf = append(sha1Buf, 0x16)
        sha1Buf, err = CalSha1(sha1Buf, r)
        if err != nil {
            return
        }
    } else { // file size > 4M
        sha1Buf = append(sha1Buf, 0x96)
        sha1BlockBuf := make([]byte, 0, blockCnt*20)
        for i := 0; i < blockCnt; i++ {
            body := io.LimitReader(r, BlockSize)
            sha1BlockBuf, err = CalSha1(sha1BlockBuf, body)
            if err != nil {
                return
            }
        }
        sha1Buf, _ = CalSha1(sha1Buf, bytes.NewReader(sha1BlockBuf))
    }
    etag = base64.URLEncoding.EncodeToString(sha1Buf)
    return
}
