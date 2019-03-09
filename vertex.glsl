#version 420 core
layout (location = 0) in vec3 vert;
layout (location = 1) in vec2 inTexCoords;

out vec3 pos;
out vec2 texCoords;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;

void main() {
	gl_Position = viewMatrix * modelMatrix * vec4(vert, 1.0);
	pos = vec3(vert);
	texCoords = inTexCoords;
}