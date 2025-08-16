package world

import (
	"path/filepath"
	"strings"
)

// RoomIDFromPath creates a stable room id from zone and file path like r01.tmj -> "zone/r01".
func RoomIDFromPath(zoneName, path string) string {
	base := filepath.Base(path)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	return zoneName + "/" + base
}