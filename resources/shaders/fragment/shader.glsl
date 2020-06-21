#version 330

uniform sampler2D tex;
uniform vec3 objectColor;
uniform vec3 lightColor;

in vec2 fragTexCoord;

out vec4 FragColor;

void main() {
    FragColor = vec4(lightColor * objectColor, 1.0);
}
