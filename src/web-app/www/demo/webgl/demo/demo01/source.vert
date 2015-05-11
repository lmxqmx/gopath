attribute vec3 aVertexPosition;
attribute vec3 aVertexColor;
uniform mat4 uModelViewMatrix;
uniform mat4 uProjectionMatrix;
varying vec4 vColor;

void main()
{
  gl_Position = uProjectionMatrix * uModelViewMatrix * vec4(aVertexPosition, 1.0);
  vColor = vec4(aVertexColor, 1.0);
}