//go:build !nosqlite
// +build !nosqlite

package main

import (
	_ "github.com/mattn/go-sqlite3"
)
