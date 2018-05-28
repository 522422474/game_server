/**
 * http相关操作

 * @author 黄承武

 * @create time 2018.4.23
 */
package web


import (
	"strings"
	// "dl_server/utils"
	// "io/ioutil"
	"strconv"
	"net/http" 
	"fmt"
	"server/datastruct"
	"server/config"
	"encoding/base64"
	// "encoding/json"
)



/**
 * 处理函数
 */
func msgHandler(w http.ResponseWriter, req *http.Request) {  
	//解决跨域问题
	w.Header().Set("Access-Control-Allow-Origin", "*");
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type");
	w.Header().Set("content-type", "application/json");

	// req.ParseForm();
	
	//读取用户发来的数据，并解析
	// msg, _ := ioutil.ReadAll(req.Body);
	data := datastruct.Netmsg{};
	// json.Unmarshal([]byte(msg), &data);

	// data.OpenId = req.FormValue("OpenId");
	// data.Id, _ = strconv.Atoi(req.FormValue("Id"));
	// data.Cmd, _ = strconv.Atoi(req.FormValue("Cmd"));
	// data.Content = req.FormValue("Content");
	// data.Sign = req.FormValue("Sign");

	datamsg, _ := base64.StdEncoding.DecodeString(req.FormValue("data"));
	datas := strings.Split(string(datamsg), "|dlgame|");
	data.OpenId = datas[0]
	data.Id, _ = strconv.Atoi(datas[1]);
	data.Cmd, _ = strconv.Atoi(datas[2]);
	data.Content = datas[3];


	fmt.Println( "接收客户端消息：", data );

	//验证消息的合法性
	// if !utils.CheckMessageSign(data.Sign, data.OpenId, data.Id, data.Cmd, data.Content) {
	// 	//与服务器算出的结果不相同，说明是不合法消息，不做任何处理，直接返回
	// 	fmt.Println("接收的消息不合法. OpenID=", data.OpenId);
	// 	return;
	// }

	// //处理消息
	DoHttpHandler(w, data);
}



/**
 * 启动http / https服务
 * @param ishttp 是否是http服务器
 */
func StartHttpServer(ishttp bool){
	
	__url := config.HTTP_URL + ":" + strconv.Itoa(config.HTTP_PORT);
	http.HandleFunc("/msgHandler", msgHandler);

	var pmsg1 int;
	if(ishttp){
		pmsg1 = 1;
	}else{ pmsg1 = 0;};
	pmsg2 := "-> 启动（主）服务... / isHttp = " + strconv.Itoa(pmsg1) + " 监听地址：" + __url;
	fmt.Println(pmsg2);
	Pushlog(pmsg2);

	//启动排行榜服务器（子线程）
	StartRankServer();

	//启动日志服务器（子线程）
	StartLogServer();

	var err error;
	if ishttp {
		//http
		err = http.ListenAndServe(__url, nil);
	}else{
		//https
		// err = http.ListenAndServeTLS(__url, nil, nil, nil); //第2个参数证书，第3个参数key.两个参数的类型：string
	}

	if err != nil {  
		pmsg3 := "-> [Error] : 启动（主）服务发生错误. " + err.Error();
		fmt.Println(pmsg3);
		Pushlog(pmsg3);
	}
}
