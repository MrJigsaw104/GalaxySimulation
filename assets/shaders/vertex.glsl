#version 410 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aVel;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform vec3 centerPos;

out vec3 particleColor;

void main() {
    gl_Position = projection * view * model * vec4(aPos, 1.0);
    
    float dist = length(aPos - centerPos);
    vec3 hotColor = vec3(1.0, 0.8, 0.0); // Sıcak sarı
    vec3 coldColor = vec3(0.0, 0.3, 1.0); // Soğuk mavi
    
    float t = clamp(dist / 300.0, 0.0, 1.0);
    particleColor = mix(hotColor, coldColor, t);
    
    gl_PointSize = mix(12.0, 6.0, t);
} 