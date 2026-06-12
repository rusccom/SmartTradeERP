package media

import "errors"

var ErrStorageNotConfigured = errors.New("media storage is not configured")
var ErrInvalidMedia = errors.New("invalid media")
