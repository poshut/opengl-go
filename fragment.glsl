#version 420 core

in vec3 pos;
in vec2 texCoords;

out vec4 color;

uniform float red;
layout(binding = 0) uniform sampler2D tex;

void main() {
	color = mix(texture(tex, texCoords), vec4(1.0,0.0,0.0,0.0), red);
	// color = vec4((pos.x+0.5), -(pos.y-0.5), 0.0, 0.0);
	// color = vec4(1.0,0.0,0.0,0.0);
}