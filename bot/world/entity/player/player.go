package player

import "github.com/Tnze/go-mc/bot/world/entity"
import "math"

// Player includes the player's status.
type Player struct {
	entity.Entity
	UUID [2]int64 //128bit UUID

	X, Y, Z    float64
	Yaw, Pitch float32
	OnGround   bool

	HeldItem int //拿着的物品栏位

	Health         float32 //血量
	Food           int32   //饱食度
	FoodSaturation float32 //食物饱和度
}

//GetPosition return the player's position
func (p *Player) GetPosition() (x, y, z float64) {
        return p.X, p.Y, p.Z
}

//GetBlockPos return the position of the Block at player's feet
func (p *Player) GetBlockPos() (x, y, z int) {
        return int(math.Floor(p.X)), int(math.Floor(p.Y)), int(math.Floor(p.Z))
}
