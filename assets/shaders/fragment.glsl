#version 410 core
in vec3 particleColor;
out vec4 FragColor;

uniform vec3 centerPos;
uniform float time;

void main() {
    vec2 center = gl_PointCoord - vec2(0.5);
    float dist = length(center);
    float alpha = smoothstep(0.5, 0.2, dist);
    
    float distToCenter = length(centerPos);
    float centralAlpha = mix(0.8, 0.2, distToCenter / 300.0);
    
    FragColor = vec4(particleColor, alpha * centralAlpha * 0.6);
}