#version 420 core
layout (location = 0) in vec3 vert;
layout (location = 1) in vec2 inTexCoords;
layout (location = 2) in vec3 normal;

out vec2 texCoords;
out vec3 toLightVector;
out vec3 surfaceNormal;

uniform mat4 modelMatrix;
uniform mat4 viewMatrix;
uniform mat4 projectionMatrix;
uniform vec3 lightPos;

void main() {
	vec4 worldPosition = modelMatrix * vec4(vert, 1.0);
	gl_Position = projectionMatrix * viewMatrix * worldPosition;
	texCoords = inTexCoords;

	// Todo: load the normal matrix as a uniform variable
	surfaceNormal = transpose(inverse(mat3(modelMatrix))) * normal;
	toLightVector = lightPos - worldPosition.xyz;
}