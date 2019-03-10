#version 420 core

in vec2 texCoords;
in vec3 surfaceNormal;
in vec3 toLightVector;

out vec4 color;

layout(binding = 0) uniform sampler2D tex;
uniform float hasTexture;

void main() {
	vec3 unitNormal = normalize(surfaceNormal);
	vec3 unitLightVector = normalize(toLightVector);
	float diffuseStrength = dot(unitLightVector, unitNormal);
	diffuseStrength = max(diffuseStrength, 0.2);

	if (hasTexture==1) {
		color = texture(tex, texCoords) * diffuseStrength;
	} else {
		color = vec4(1.0,1.0,1.0,0.0) * diffuseStrength;
	}
}