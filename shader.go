package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// ShaderProgram represents a shader program containing a vertex and a fragment shader
type ShaderProgram uint32

// Delete deletes the OpenGL shader program
func (program ShaderProgram) Delete() {
	gl.DeleteProgram(uint32(program))
}

// Use makes the shader program the active one
func (program ShaderProgram) Use() {
	gl.UseProgram(uint32(program))
}

// Unuse makes no shader program the active one
func (program ShaderProgram) Unuse() {
	gl.UseProgram(0)
}

// CreateProgramFromFiles creates a shader program from the vertex and fragment shader paths
func CreateProgramFromFiles(vertex string, fragment string) (ShaderProgram, error) {
	vertexShader, err := readShaderFile(vertex)
	if err != nil {
		return 0, err
	}
	fragmentShader, err := readShaderFile(fragment)
	if err != nil {
		return 0, err
	}

	return CreateProgramFromSource(vertexShader, fragmentShader)
}

// CreateProgramFromSource creates a shader program from shader sources
func CreateProgramFromSource(vertex string, fragment string) (ShaderProgram, error) {
	vertexShader, err := compileShader(vertex, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	fragmentShader, err := compileShader(fragment, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
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
		return 0, errors.New("Couldn't link program: \n" + infoLog)
	}
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return ShaderProgram(shader), nil
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
