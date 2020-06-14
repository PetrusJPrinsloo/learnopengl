#version 330 core

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main()
{
    fragTexCoord = vertTexCoord;
    gl_Position = projection * view * model * vec4(vert, 1);
}
