#version 420 core

in vec3 pos;
in vec2 texCoords;

out vec4 color;

layout(binding = 0) uniform sampler2D tex;

void main() {
	color = texture(tex, texCoords);
}