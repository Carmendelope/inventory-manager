/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestInventoryPackage(t *testing.T){
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Inventory Handler & Manager package suite")
}