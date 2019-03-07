#version 400 core
in vec3 pos;
out vec4 color;
void main() {
	color = vec4((pos.x+0.5), -(pos.y-0.5), 0.0, 0.0);
}