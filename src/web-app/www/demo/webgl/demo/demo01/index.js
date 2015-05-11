/**
 * Created by L on 2015/5/11.
 */


window.onload = function(){
    var webgl = document.getElementById('webgl');
    var gl = webgl.getContext('webgl');


    Promise.all([get('source.vert'), get('source.frag')]).then(function(sources){
        var vertexSource = sources[0];
        var fragmentSource = sources[1];


        // 编译顶点着色器
        var vertexShader = gl.createShader(gl.VERTEX_SHADER);
        gl.shaderSource(vertexShader, vertexSource);
        gl.compileShader(vertexShader);

        // 编译片段着色器
        var fragmentShader = gl.createShader(gl.FRAGMENT_SHADER);
        gl.shaderSource(fragmentShader, fragmentSource);
        gl.compileShader(fragmentShader);

        // 链接到着色器程序
        var program = gl.createProgram();
        gl.attachShader(program, vertexShader);
        gl.attachShader(program, fragmentShader);
        gl.linkProgram(program);
        gl.useProgram(program);

        var modelViewMatrix = mat4.create();//模型视图矩阵
        mat4.lookAt(
            modelViewMatrix,
            [4, 4, 8], //观察者站在(4,4,8)坐标位置
            [0, 0, 0], //眼睛望向(0,0,0)坐标位置
            [0, 1, 0] //头顶朝向(0,1,0)坐标位置
        );

        //获得变量在内存里的地址，再在相应的内存里写入值。
        var uModelViewMatrix = gl.getUniformLocation(program, "uModelViewMatrix");
        gl.uniformMatrix4fv(
            uModelViewMatrix,
            false, //false参数表示不需要转置（行变列，列变行）这个矩阵
            modelViewMatrix);

        //设置投影矩阵
        var projectionMatrix = mat4.create();
        mat4.perspective(
            projectionMatrix,
            Math.PI / 6, //视角为30°
            webgl.width / webgl.height,//视口的宽高比
            0.1,//视锥近截面到观察点的距离
            100//视锥远截面到观察点的距离
        );

        // 传递给顶点着色器的uniform变量uProjectionMatrix
        var uProjectionMatrix = gl.getUniformLocation(program, "uProjectionMatrix");
        gl.uniformMatrix4fv(uProjectionMatrix, false, projectionMatrix);

        var vertices = [
            //前
            1.0, 1.0, 1.0, 0.0, 0.8, 0.0,//前三个元素表示该顶点的坐标XYZ，后三个元素表示该顶点的颜色RGB
            -1.0, 1.0, 1.0, 0.0, 0.8, 0.0,
            -1.0, -1.0, 1.0, 0.0, 0.8, 0.0,
            1.0, -1.0, 1.0, 0.0, 0.8, 0.0,
            //后
            1.0, 1.0, -1.0, 0.6, 0.9, 0.0,
            -1.0, 1.0, -1.0, 0.6, 0.9, 0.0,
            -1.0, -1.0, -1.0, 0.6, 0.9, 0.0,
            1.0, -1.0, -1.0, 0.6, 0.9, 0.0,
            //上
            1.0, 1.0, -1.0, 1.0, 1.0, 0.0,
            -1.0, 1.0, -1.0, 1.0, 1.0, 0.0,
            -1.0, 1.0, 1.0, 1.0, 1.0, 0.0,
            1.0, 1.0, 1.0, 1.0, 1.0, 0.0,
            //下
            1.0, -1.0, -1.0, 1.0, 0.5, 0.0,
            -1.0, -1.0, -1.0, 1.0, 0.5, 0.0,
            -1.0, -1.0, 1.0, 1.0, 0.5, 0.0,
            1.0, -1.0, 1.0, 1.0, 0.5, 0.0,
            //右
            1.0, 1.0, -1.0, 0.9, 0.0, 0.2,
            1.0, 1.0, 1.0, 0.9, 0.0, 0.2,
            1.0, -1.0, 1.0, 0.9, 0.0, 0.2,
            1.0, -1.0, -1.0, 0.9, 0.0, 0.2,
            //左
            -1.0, 1.0, -1.0, 0.6, 0.0, 0.6,
            -1.0, 1.0, 1.0, 0.6, 0.0, 0.6,
            -1.0, -1.0, 1.0, 0.6, 0.0, 0.6,
            -1.0, -1.0, -1.0, 0.6, 0.0, 0.6
        ];

        // vertices只是一个数组，WebGL并不能直接操作JS数组，
        // 我们需要把它转换成类型化数组然后载入缓冲区
        var vertexBuffer = gl.createBuffer();
        gl.bindBuffer(gl.ARRAY_BUFFER, vertexBuffer);
        gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(vertices),
            gl.STATIC_DRAW //数据只加载一次
        );

        var aVertexPosition = gl.getAttribLocation(program, "aVertexPosition");
        gl.vertexAttribPointer(aVertexPosition, 3, gl.FLOAT, false, 24, 0);
        gl.enableVertexAttribArray(aVertexPosition);

        var aVertexColor = gl.getAttribLocation(program, "aVertexColor");
        gl.vertexAttribPointer(aVertexColor, 3, gl.FLOAT, false, 24, 12);
        gl.enableVertexAttribArray(aVertexColor);

        //绘制顺序的信息
        var indices = [
            0, 1, 2, 0, 2, 3,
            4, 6, 5, 4, 7, 6,
            8, 9, 10, 8, 10, 11,
            12, 14, 13, 12, 15, 14,
            16, 17, 18, 16, 18, 19,
            20, 22, 21, 20, 23, 22
        ];

        var indexBuffer = gl.createBuffer();
        gl.bindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer);
        gl.bufferData(gl.ELEMENT_ARRAY_BUFFER, new Uint8Array(indices), gl.STATIC_DRAW);

        //绘制
        gl.enable(gl.DEPTH_TEST);
        gl.enable(gl.CULL_FACE);
        gl.clearColor(0.0, 0.0, 0.0, 1.0);
        gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
        gl.drawElements(gl.TRIANGLES, indices.length, gl.UNSIGNED_BYTE, 0);
    });


};

function get(url){
    return new Promise(function(resolve){
        var xhr = new XMLHttpRequest();
        xhr.onload = function(){
            resolve(this.responseText);
        };
        xhr.open("get", url);
        xhr.send();
    });
}

