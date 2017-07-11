/**
 * Created by Administrator on 2017/7/10.
 */
function SSE(url,func,retryTime) {
    var es=null;
    var keep_aliveSecs=(retryTime)?parseInt(retryTime):0;//重连的时间间隔
    var keep_aliveTimer=null;//重置倒计时

    connect();
    
    //连接函数
    function connect() {
        if (es)es.close();
        if(keep_aliveSecs!=0){
            gotActivity();
        }
        if (window.EventSource){
            startEventSource();
        }else {
            console.log("Your browser does not support SSE");
        }
    }

    //重连函数
    function gotActivity() {
        if(keep_aliveTimer!=null)
            clearTimeout(keep_aliveSecs);
        keep_aliveTimer=setTimeout(connect,keep_aliveSecs * 1000);
    }
    
    //开启SSE函数
    function startEventSource() {
        es=new EventSource(url);
        es.addEventListener("message", (e) => {
            processOneLine(e.data);
        },false);
        es.addEventListener("error",handleError,false);
    }
    
    //处理数据推送的函数
    function  processOneLine(s) {
        try {
            var d=JSON.parse(s);
            func(d);
        }catch (e){
            console.log("BAD JSON:"+s+"\n"+e);
        }
    }
    
    //错误处理函数
    function handleError(e) {
        switch (e.target.readyState){
            case 0:
                msg="connecting";
                break;
            case 1:
                msg="connect successful";
                break;
            case 2:
                msg="the url has problem,please check the url";
                break;
            default:
                msg="error is unknown"
        }
        console.log("BAD Error:"+msg);
    }
}
