/*
Copyright © 2021-2023 Macaroni OS Linux
See AUTHORS and LICENSE for the license details and contributors.
*/
package reposcan_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSolver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reposcan Suite")
}
