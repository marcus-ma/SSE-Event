<?php
header("Content-Type:text/event-stream");//数据推送协议的头函数格式
header('Cache-Control: no-cache,must-revalidate');//阻止缓存(HTTP1.1规范)
header('Expires: Sun, 31 Dec 2000 05:00:00 GMT');//（旧浏览器规范的阻止缓存）
set_time_limit(0);//阻止死亡（此脚本挂掉30秒后，可进行修复）

function sendData($data){
    
    echo "data:";//数据推送协议的开头格式
    echo json_encode($data)."\n";
    echo "\n";//数据推送协议的结尾格式
    
    //不把数据缓存，立即将数据进行推送，
    ob_flush();
    flush();
}

while (true){
    $time = date("Y-m-d H:i:s");
    $data=array('msg'=>'The server time is:'.$time);
    
    sendData($data);
    //停顿秒速
    sleep(1);
}





