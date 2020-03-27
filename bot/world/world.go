package world

import (
	"github.com/Tnze/go-mc/bot/world/entity"
	"github.com/Tnze/go-mc/data"
)

//World record all of the things in the world where player at
type World struct {
	Entities map[int32]entity.Entity
	Chunks   map[ChunkLoc]*Chunk
}

//Chunk store a 256*16*16 clolumn blocks
type Chunk struct {
	sections [16]Section
}

//Section store a 16*16*16 cube blocks
type Section struct {
	blocks [16][16][16]Block
}

//Block is the base of world
type Block struct {
	id uint
}

type ChunkLoc struct {
	X, Z int
}

// //Entity 表示一个实体
// type Entity interface {
// 	EntityID() int32
// }

// //Face is a face of a block
// type Face byte

// // All six faces in a block
// const (
// 	Bottom Face = iota
// 	Top
// 	North
// 	South
// 	West
// 	East
// )

// getBlock return the block in the position (x, y, z)
func (w *World) GetBlock(x, y, z int) Block {
	c := w.Chunks[ChunkLoc{x >> 4, z >> 4}]
	if c != nil {
		cx, cy, cz := x&15, y&15, z&15
		/*
			n = n&(16-1)

			is equal to

			n %= 16
			if n < 0 { n += 16 }
		*/

		return c.sections[y/16].blocks[cx][cy][cz]
	}

	return Block{id: 0}
}

func (b Block) String() string {
	return data.BlockNameByID[b.id]
}

func (b Block) ID() uint {
	return b.id
}

//LoadChunk load chunk at (x, z)
func (w *World) LoadChunk(x, z int, c *Chunk) {
	w.Chunks[ChunkLoc{X: x, Z: z}] = c
}
