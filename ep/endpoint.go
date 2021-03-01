package ep

import (
	"github.com/harriklein/pBE/pBEServer/ep/auth"
	"github.com/harriklein/pBE/pBEServer/ep/docs"
	"github.com/harriklein/pBE/pBEServer/ep/dummy"
)

// EndPointsInit initializes all endpoints
func Init() {
	auth.Init()
	dummy.Init()
	docs.Init()
}
