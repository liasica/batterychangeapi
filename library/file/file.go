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
    "github.com/gabriel-vasile/mimetype"
    "mime/multipart"
)

type MultipartFile struct {
    *multipart.FileHeader
}

func (f *MultipartFile) GetFileContentType() (*mimetype.MIME, error) {
    r, err := f.Open()

    if err != nil {
        return nil, err
    }
    defer r.Close()

    return mimetype.DetectReader(r)
}
