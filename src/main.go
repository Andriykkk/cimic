package main

import (
	"runtime"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

var (
	cameraPos   = mgl32.Vec3{0.0, 0.0, 3.0}
	cameraFront = mgl32.Vec3{0.0, 0.0, -1.0}
	cameraUp    = mgl32.Vec3{0.0, 1.0, 0.0}

	cameraSpeed = float32(250.0)

	nearClipping = 0.1
	farClipping  = 10000.0
	

	deltaTime = float32(0.0)
	lastFrame = float32(0.0)

	lastX, lastY = float64(windowWidth) / 2, float64(windowHeight) / 2
	firstMouse   = true
	yaw          = -90.0
	pitch        = 0.0
	sensitivity  = 0.1
)

type Vertex struct {
	X, Y, Z float64
}

type Face struct {
	VertexIndices  []int
	TextureIndices []int
	NormalIndices  []int
}

func init() {
	runtime.LockOSThread()
}

func main() {
	vertices, faces, _, _, _ := readObj("plant.obj")
	initWindow(vertices, faces)
}
