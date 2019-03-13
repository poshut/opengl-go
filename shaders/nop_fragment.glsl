#version 420 core

in vec2 pass_texCoords;

out vec4 color;

layout (binding = 0) uniform sampler2D tex;

void main() {
    color = texture(tex, pass_texCoords);
    // Invert the colors:
    // color = vec4(1 - color.xyz, 1.0);
}