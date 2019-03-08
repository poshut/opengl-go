package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// ShaderProgram represents a shader program containing a vertex and a fragment shader
type ShaderProgram struct {
	programID        uint32
	uniformLocations map[string]int32
}

// Delete deletes the OpenGL shader program
func (program ShaderProgram) Delete() {
	gl.DeleteProgram(uint32(program.programID))
}

// Use makes the shader program the active one
func (program ShaderProgram) Use() {
	gl.UseProgram(uint32(program.programID))
}

// Unuse makes no shader program the active one
func (program ShaderProgram) Unuse() {
	gl.UseProgram(0)
}

// LoadUniformFloat loads the given value to the given location. THE PROGRAM MUST BE ACTIVE!
func (program *ShaderProgram) LoadUniformFloat(location string, variable float32) {
	gl.Uniform1f(program.GetUniformLocation(location), variable)
}

// LoadUniformVector loads the given value to the given location. THE PROGRAM MUST BE ACTIVE!
func (program *ShaderProgram) LoadUniformVector(location string, variable mgl32.Vec3) {
	gl.Uniform3f(program.GetUniformLocation(location), variable.X(), variable.Y(), variable.Z())

}

// LoadUniformMatrix loads the given value to the given location. THE PROGRAM MUST BE ACTIVE!
func (program *ShaderProgram) LoadUniformMatrix(location string, variable mgl32.Mat4) {
	gl.UniformMatrix4fv(program.GetUniformLocation(location), 1, false, &variable[0])
}

// GetUniformLocation returns the uniform variable location of the variable identifier. THE PROGRAM MUST BE ACTIVE!
func (program *ShaderProgram) GetUniformLocation(s string) int32 {
	if i, ok := program.uniformLocations[s]; ok {
		return i
	}
	location := gl.GetUniformLocation(program.programID, gl.Str(s+"\x00"))
	program.uniformLocations[s] = location
	return location
}

// CreateProgramFromFiles creates a shader program from the vertex and fragment shader paths
func CreateProgramFromFiles(vertex string, fragment string) (ShaderProgram, error) {
	vertexShader, err := readShaderFile(vertex)
	if err != nil {
		return ShaderProgram{0, nil}, err
	}
	fragmentShader, err := readShaderFile(fragment)
	if err != nil {
		return ShaderProgram{0, nil}, err
	}

	return CreateProgramFromSource(vertexShader, fragmentShader)
}

// CreateProgramFromSource creates a shader program from shader sources
func CreateProgramFromSource(vertex string, fragment string) (ShaderProgram, error) {
	vertexShader, err := compileShader(vertex, gl.VERTEX_SHADER)
	if err != nil {
		return ShaderProgram{0, nil}, err
	}
	fragmentShader, err := compileShader(fragment, gl.FRAGMENT_SHADER)
	if err != nil {
		return ShaderProgram{0, nil}, err
	}

	shader := gl.CreateProgram()
	gl.AttachShader(shader, vertexShader)
	gl.AttachShader(shader, fragmentShader)
	gl.LinkProgram(shader)

	var status int32
	gl.GetProgramiv(shader, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shader, gl.INFO_LOG_LENGTH, &status)
		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shader, logLength, nil, gl.Str(infoLog))
		return ShaderProgram{0, nil}, errors.New("Couldn't link program: \n" + infoLog)
	}
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return ShaderProgram{shader, make(map[string]int32)}, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	cstr, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, cstr, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(infoLog))
		return 0, errors.New("Couldn't compile shader: \n" + infoLog)
	}
	return shader, nil
}

func readShaderFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data) + "\x00", nil
}
