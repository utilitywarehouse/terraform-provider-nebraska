package nebraska

import (
	"errors"
)

// Arch represents an arch supported by Nebraska
type Arch uint

const (

	// ArchAll = all
	ArchAll Arch = iota
	// ArchAMD64 = amd64
	ArchAMD64
	// ArchAArch64 = aarch64
	ArchAArch64
	// ArchX86 = x86
	ArchX86
)

var (
	// ErrInvalidArch is a custom error returned when an unsupported arch is
	// requested
	ErrInvalidArch = errors.New("nebraska: invalid/unsupported arch")

	// ValidArchs are the archs that Nebraska supports
	ValidArchs = []string{
		"all",
		"amd64",
		"aarch64",
		"x86",
	}
)

// String returns the string representation of the arch
func (a Arch) String() string {
	i := int(a)

	return ValidArchs[i]
}

// ArchFromString parses the string into an Arch
func ArchFromString(s string) (Arch, error) {
	for i, sd := range ValidArchs {
		if s == sd {
			return Arch(i), nil
		}

	}

	return ArchAll, ErrInvalidArch
}
