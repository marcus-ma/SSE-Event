SSE
  全称Server-Sent Event(服务端推送事件)，一项基于HTML5的允许服务端向客户端推送新数据的技术。

使用方法

    前端：
（1）在要使用SSE功能的HTML5文件中引入封装好的SSE.js文件

    <script src="SSE.js"></script>

（2）在js的部分先声明几个变量：
        第一个变量（必要）：要连接的URL，
        第二个变量（必要）：处理函数，作用的处理由服务端推送过来的JSON数据
        第三个变量（选要）：设置重连的时间，以秒作为标准单位，当服务端故障可以重新连接
    变量设置完后，创建SSE实例，将变量传到参数中去

    <script>
        //此处连接SSE.php
         url="SSE.php";

         //处理函数，函数的参数为已经解析了的JSON数据
         func=function (d) {
            //此处为要做出的处理
             document.getElementById("x").innerHTML +="\n"+d.msg;
         };

         //声明重连的秒数变量，即几秒后再进行接收推送
         retry_time=2;

        //设置倒计时，创建SSE实例，将变量传到参数中去
        setTimeout(function () {
            new SSE(url,func,retry_time);
        },100);
    </script>


  后端
（3）以PHP为例，服务端的SSE协议规格在SSE.php文件中已经做了详细说明，可以打开该文件进行参考


