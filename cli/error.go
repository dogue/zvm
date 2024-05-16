// Copyright 2022 Tristan Isham. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cli

import (
	"errors"
)

var (
	ErrMissingBundlePath     = errors.New("bundle download path not found")
	ErrUnsupportedSystem     = errors.New("unsupported system for Zig")
	ErrUnsupportedVersion    = errors.New("unsupported Zig version")
	ErrMissingInstallPathEnv = errors.New("env 'ZVM_INSTALL' is not set")
	ErrFailedUpgrade         = errors.New("failed to self-upgrade zvm")
	ErrInvalidVersionMap     = errors.New("invalid version map format")
	ErrInvalidZlsVersion     = errors.New("invalid ZLS version")
)
