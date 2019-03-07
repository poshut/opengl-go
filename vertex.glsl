#version 400 core
layout (location = 0) in vec3 vert;
out vec3 pos;
void main() {
	gl_Position = vec4(vert, 1.0);
	pos = vec3(vert);
}