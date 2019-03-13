package main

import (
	"fmt"
	"runtime"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const windowWidth = 1080
const windowHeight = 720

func main() {
	// Create a window and initialize OpenGL
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
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "random window", nil, nil)
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

	// Load model and texture and create an entity
	model, err := CreateModelFromFile("res/dragon.obj")
	if err != nil {
		panic(err)
	}
	// err = model.AddTexture("res/stallTexture.png", true)
	// if err != nil {
	// 	panic(err)
	// }
	defer model.Delete()

	entity := Entity{mgl32.Vec3{0.0, -5.0, -20.0}, mgl32.Vec3{0.0, 0.0, 0.0}, 1.0, &model}

	// Load the shader
	program, err := CreateProgramFromFiles("shaders/vertex.glsl", "shaders/fragment.glsl")
	if err != nil {
		panic(err)
	}
	defer program.Delete()

	fboProgram, err := CreateProgramFromFiles("shaders/nop_vertex.glsl", "shaders/nop_fragment.glsl")
	if err != nil {
		panic(err)
	}
	defer program.Delete()

	// Load the perspective matrix
	projectionMatrix := mgl32.Perspective(mgl32.DegToRad(fov), float32(windowWidth)/float32(windowHeight), nearPlane, farPlane)
	program.Use()
	program.LoadUniformMatrix("projectionMatrix", projectionMatrix)
	program.Unuse()

	// Create the camera
	camera := NewCamera(window)

	InitializeFramebuffers()
	fbo, err := NewFramebuffer()
	if err != nil {
		panic(err)
	}
	err = fbo.AddColorAttachment()
	if err != nil {
		panic(err)
	}
	fbo.AddRenderbufferDepthAndStencil()
	if !fbo.IsComplete() {
		fmt.Println(gl.CheckFramebufferStatus(fbo.id))
		panic("fbo not complete")
	}
	defer fbo.Delete()

	// Enable depth testing
	gl.Enable(gl.DEPTH_TEST)
	glfw.SwapInterval(1)

	for !window.ShouldClose() {

		// Update status
		camera.Update(window)
		entity.rotation = entity.rotation.Add(mgl32.Vec3{0.0, 0.01, 0.0})

		gl.Enable(gl.DEPTH_TEST)
		gl.ClearColor(0.0, 0.0, 0.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render scene to framebuffer
		fbo.Use()
		program.Use()
		program.LoadUniformVector("lightPos", mgl32.Vec3{0.0, 0.0, 0.0})
		camera.Load(&program)
		model.Bind(&program)
		entity.Load(&program)
		model.Draw()
		model.Unbind(&program)
		program.Unuse()
		fbo.Unuse()

		// Render framebuffer to screen
		fboProgram.Use()
		gl.Disable(gl.DEPTH_TEST)
		fbo.Draw(&fboProgram)
		fboProgram.Unuse()

		window.SwapBuffers()
		glfw.PollEvents()

		e := gl.GetError()
		if e != gl.NO_ERROR {
			fmt.Println("a gl error has occured")
		}
	}
}
