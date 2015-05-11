var gl;

function webGLStart(){
    var canvas = document.getElementById('lesson01-canvas');
    initGL(canvas);
}

function initGL(canvas) {
    try {
        gl = canvas.getContext("experimental-webgl");
        gl.viewportWidth = canvas.width;
        gl.viewportHeight = canvas.height;
    } catch(e) {
        console.log(e);
    }
    if (!gl) {
        alert("Could not initialise WebGL, sorry :-( ");
    }
}