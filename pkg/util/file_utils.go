// Copyright (c) 2021 Yandex LLC. All rights reserved.
// Author: Martynov Pavel <covariance@yandex-team.ru>

package util

import "os"

func Exists(path string) (os.FileInfo, error) {
	res, err := os.Stat(path)
	if err == nil {
		return res, nil
	}
	if os.IsNotExist(err) {
		return nil, nil
	}
	return nil, err
}
