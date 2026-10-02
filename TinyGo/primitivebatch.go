package main


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

func (pb *PrimitiveBatch) End(g *Game) {
	pb.Flush(g)
	pb.hasBegun = false
}

func (pb *PrimitiveBatch) AddVertex(vertex Vector2, g *Game) {
	var newPrimitive = (pb.positionInBuffer % pb.numVertsPerPrimitive) == 0

	if newPrimitive &&
		(pb.positionInBuffer+pb.numVertsPerPrimitive) >= len(pb.vertices) {
			pb.Flush(g)
	}

	pb.vertices[pb.positionInBuffer] = vertex
	pb.positionInBuffer++
}

func (pb *PrimitiveBatch) Flush(g *Game) {

	if pb.positionInBuffer == 0 {
		return
	}

	g.ctx.Set("fillStyle", "#1504fa")
	for i := 0; i < len(pb.vertices)-2; i += 2 {
		g.ctx.Call("fillRect", float32(pb.vertices[i].x),
			float32(pb.vertices[i].y), 20, 280)
	}

	pb.positionInBuffer = 0
}
