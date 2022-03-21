package nebraska

import "errors"

var (
	// ErrInvalidArch is a custom error returned when an unsupported arch is
	// requested
	ErrInvalidArch = errors.New("nebraska: invalid/unsupported arch")

	// ValidArchs are the archs that Nebraska supports
	// https://github.com/kinvolk/nebraska/blob/main/backend/pkg/api/arch.go#L37-L43
	ValidArchs = []string{
		"all",
		"amd64",
		"aarch64",
		"x86",
	}
)
