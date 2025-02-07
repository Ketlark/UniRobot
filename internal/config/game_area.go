package config

const (
	GLOBAL_BUFFER_SIZE    = 100 * 1024 // 100KB for global area
	POSITION_BUFFER_SIZE  = 10 * 1024  // 10KB for position area
	FIGHT_BUFFER_SIZE     = 20 * 1024  // 20KB for fight state area
	INVENTORY_BUFFER_SIZE = 30 * 1024  // 30KB for inventory area
)

const (
	SWIFT_CAPTURE_LIBRARY_PATH = "./bin/ScreenshotHelper.dylib"
)

type Rectangle struct {
	X, Y, Width, Height int
}

type GameArea struct {
	Name      string
	Size      uint32
	Rectangle Rectangle
}

// GameAreas holds all game areas in a map for O(1) lookup
var GameAreas = map[string]GameArea{
	"position": {
		Name:      "position",
		Size:      POSITION_BUFFER_SIZE,
		Rectangle: Rectangle{8, 160, 128, 45},
	},
	"fight": {
		Name:      "fight",
		Size:      FIGHT_BUFFER_SIZE,
		Rectangle: Rectangle{800, 600, 200, 100},
	},
}

func GetGameAreas() []GameArea {
	areas := make([]GameArea, 0, len(GameAreas))
	for _, area := range GameAreas {
		areas = append(areas, area)
	}
	return areas
}

func GetGameAreaByName(name string) (GameArea, bool) {
	area, exists := GameAreas[name]
	return area, exists
}
