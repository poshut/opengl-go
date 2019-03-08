package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	runtime.LockOSThread()
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	window, err := glfw.CreateWindow(640, 480, "random window", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.Enable(gl.DEBUG_OUTPUT)

	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,
		0.5, 0.5, 0.0,
	}

	indices := []uint32{
		0, 1, 2,
		1, 2, 3,
	}

	texCoords := []float32{
		0.0, 0.0,
		1.0, 0.0,
		0.0, 1.0,
		1.0, 1.0,
	}

	model, err := CreateModelFromData(vertices, indices, texCoords)
	// model, err := CreateModelFromFile("bunny.obj")
	if err != nil {
		panic(err)
	}
	err = model.AddTexture("brick.jpeg", true)
	if err != nil {
		panic(err)
	}
	defer model.Delete()

	program, err := CreateProgramFromFiles("vertex.glsl", "fragment.glsl")
	if err != nil {
		panic(err)
	}
	defer program.Delete()

	var red float32

	for !window.ShouldClose() {
		gl.ClearColor(1.0, 1.0, 1.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		program.Use()
		program.LoadUniformFloat("red", red)
		red += 0.01
		if red > 1 {
			red = 0.0
		}
		model.Draw()
		program.Unuse()

		window.SwapBuffers()
		glfw.PollEvents()

		e := gl.GetError()
		if e != gl.NO_ERROR {
			fmt.Println("a gl error has occured")
		}
	}
}
