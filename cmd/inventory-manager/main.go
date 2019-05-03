/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package main

import (
	"github.com/nalej/inventory-manager/cmd/inventory-manager/commands"
	"github.com/nalej/inventory-manager/version"
)

// MainVersion with the application version.
var MainVersion string
// MainCommit with the commit id.
var MainCommit string

func main() {
	version.AppVersion = MainVersion
	version.Commit = MainCommit
	commands.Execute()
}
