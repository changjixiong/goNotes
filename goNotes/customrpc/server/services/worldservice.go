package services

type WorldService struct {
}

var DefaultWorldService = &WorldService{}

func (worldService *WorldService) World(num int) string {

	// return strings.Repeat("World", num)
	return "World"
}
