/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package edgecontroller

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestEdgeControllerPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "EdgeController Handler & Manager package suite")
}
