### First Rooms + Movement Tuning

This milestone connects the initial rooms and dials in baseline player feel.

- Connected rooms: `main` → `forest_right` → `forest_left` (and back)
- Transitions and spawns are data-driven via Tiled TMJ data (`portals` and `spawns` layers)
- Spawns can set initial player facing using `facing_id` ("east"/"west"/"right"/"left")
- Movement tuning adjustments: higher jump power, stronger air control, corrected coyote time

#### How to add a room

1) Create a Tiled map in `data/zones/<zone>/` (see existing zone files for reference).
2) Define portals and optional spawn points within the TMJ file so the loader can handle room creation.

Example spawn with facing:

```json
{"id": "west", "x": 32, "y": 96, "facing_id": "east"}
```

#### Movement knobs (defaults)

Configured in `engine/config.go` under `PlayerPhysics`:
- MoveSpeed: 2
- JumpPower: 8
- AirControl: 0.85
- Friction: 1
- AirFriction: 0
- CoyoteTime: 6
- JumpBufferTime: 10
- VariableJumpHeight: true
- MinJumpHeight: 0.5
- Gravity: 1
- MaxFallSpeed: 12
- FastFallMultiplier: 1.75

You can also tweak during runtime via Settings → Developer tab where applicable.

#### Known TODOs

- Minimap rendering is scaffolded in `world/minimap.go` (data API ready, rendering TBD)
- Room display names still mirror `zone_id` until the `Room` interface is extended