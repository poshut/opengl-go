#version 420 core
in vec3 pos;
in vec2 texCoords;
out vec4 color;
layout(binding = 0) uniform sampler2D tex;
void main() {
	color = texture(tex, texCoords);
	// color = vec4((pos.x+0.5), -(pos.y-0.5), 0.0, 0.0);
	// color = vec4(1.0,0.0,0.0,0.0);
}