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

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	defer gl.DeleteVertexArrays(1, &vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	defer gl.DeleteBuffers(1, &vbo)

	program, err := CreateProgramFromFiles("vertex.glsl", "fragment.glsl")
	if err != nil {
		panic(err)
	}
	// gl.BindAttribLocation(program, 0, gl.Str("vert"+"\x00"))
	defer program.Delete()

	for !window.ShouldClose() {
		gl.ClearColor(1.0, 1.0, 1.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		program.Use()
		gl.BindVertexArray(vao)
		gl.EnableVertexAttribArray(0)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, int32(len(vertices)))
		gl.DisableVertexAttribArray(0)
		gl.BindVertexArray(0)
		program.Unuse()

		window.SwapBuffers()
		glfw.PollEvents()

		e := gl.GetError()
		if e != gl.NO_ERROR {
			fmt.Println("a gl error has occured")
		}
	}
}
