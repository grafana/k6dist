package k6dist

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidPlatform thrown when platform cannot be parsed.
var ErrInvalidPlatform = errors.New("invalid platform")

// Platform contains a target platform (OS and architecture) for building.
type Platform struct {
	OS   string
	Arch string
}

// String returns the platform in string format.
func (p Platform) String() string {
	return p.OS + "/" + p.Arch
}

// ParsePlatform parses string representation of Platform.
func ParsePlatform(value string) (*Platform, error) {
	idx := strings.IndexRune(value, '/')
	if idx <= 0 || idx == len(value)-1 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidPlatform, value)
	}

	return &Platform{OS: value[:idx], Arch: value[idx+1:]}, nil
}
