//go:build tools
// +build tools

package tools

// nolint
import (
	_ "github.com/onsi/ginkgo/ginkgo"
	_ "github.com/onsi/ginkgo/v2/ginkgo/generators"
	_ "golang.org/x/tools/cmd/goimports"
)
