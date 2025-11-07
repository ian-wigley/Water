package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PrimitiveBatch struct {
	defaultBufferSize    int
	vertices             []Vector2
	positionInBuffer     int
	numVertsPerPrimitive int
	hasBegun             bool
}

func (pb *PrimitiveBatch) Construct() {
	pb.defaultBufferSize = 500
	pb.vertices = make([]Vector2, 0, pb.defaultBufferSize)
	for i := 0; i < pb.defaultBufferSize; i++ {
		var v Vector2
		pb.vertices = append(pb.vertices, v)
	}
	pb.positionInBuffer = 0
	pb.numVertsPerPrimitive = 2
	pb.hasBegun = false
}

func (pb *PrimitiveBatch) Begin() {
	pb.hasBegun = true
}

func (pb *PrimitiveBatch) End(screen *ebiten.Image) {
	pb.Flush(screen)
	pb.hasBegun = false
}

func (pb *PrimitiveBatch) AddVertex(vertex Vector2, screen *ebiten.Image) {

	var newPrimitive = (pb.positionInBuffer % pb.numVertsPerPrimitive) == 0

	if newPrimitive &&
		(pb.positionInBuffer+pb.numVertsPerPrimitive) >= len(pb.vertices) {
		pb.Flush(screen)
	}

	pb.vertices[pb.positionInBuffer] = vertex
	pb.positionInBuffer++
}

func (pb *PrimitiveBatch) Flush(screen *ebiten.Image) {

	if pb.positionInBuffer == 0 {
		return
	}

	var colour = color.RGBA{R: 0, G: 0, B: 255, A: 250}
	for i := 0; i < len(pb.vertices)-2; i += 2 {
		vector.FillRect(screen, float32(pb.vertices[i].x),
			float32(pb.vertices[i].y), 20, 280, colour, false)
	}

	pb.positionInBuffer = 0
}
