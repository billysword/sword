package world

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sword/engine"
	"sword/internal/tiled"
)

// LoadZoneRoomsFromData scans data/zones/<zoneName> for .tmj files and registers rooms and transitions
func LoadZoneRoomsFromData(rtm *RoomTransitionManager, zoneName string, baseDir string) error {
	zoneDir := filepath.Join(baseDir, "data", "zones", zoneName)
	entries, err := os.ReadDir(zoneDir)
	if err != nil {
		return fmt.Errorf("read zone dir: %w", err)
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".tmj") {
			continue
		}
		path := filepath.Join(zoneDir, e.Name())
		lm, err := tiled.LoadMap(path)
		if err != nil {
			engine.LogInfo(fmt.Sprintf("failed to load tmj %s: %v", path, err))
			continue
		}
		roomID := RoomIDFromPath(zoneName, e.Name())
		room := NewTiledRoomFromLoadedMap(roomID, lm)
		rtm.RegisterRoom(room)

		// Add transitions from portals
		u := engine.GetPhysicsUnit()
		for _, p := range lm.Portals {
			toZone := p.ToZone
			toRoom := p.ToRoom
			toPortal := p.ToPortal
			if toZone == "" {
				toZone = zoneName
			}
			if toRoom == "" {
				continue
			}
			targetID := toZone + "/" + toRoom
			// Trigger rect in physics units
			rx := int(p.RectPx[0])
			ry := int(p.RectPx[1])
			rw := int(p.RectPx[2])
			rh := int(p.RectPx[3])
			// Convert from pixel coordinates in TMJ to our physics units (tile px)
			// TMJ uses same base pixels as physics unit; numbers are already in px
			trigger := Rectangle{X: rx, Y: ry, Width: rw, Height: rh}
			// Direction heuristic based on portal name
			dir := directionFromPortalName(p.Name)
			transition := TransitionPoint{
				Type:          TransitionDoor,
				TriggerBounds: trigger,
				TargetRoomID:  targetID,
				TargetSpawnID: toPortal,
				IsEnabled:     true,
				Direction:     dir,
			}
			_ = u // reserved; may use scaling if needed
			if err := rtm.AddTransitionPoint(roomID, transition); err != nil {
				engine.LogInfo("failed adding transition: " + err.Error())
			}
		}
	}
	return nil
}

func directionFromPortalName(name string) Direction {
	s := strings.ToLower(name)
	switch s {
	case "left":
		return West
	case "right":
		return East
	case "up", "top":
		return North
	case "down", "bottom":
		return South
	default:
		return East
	}
}