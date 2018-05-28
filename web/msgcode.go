/**
 * 用于声明和定义前、后端通讯消息码
 * 消息码前缀定义：
		1. 接收客户端发送的消息码前缀：CS_
		2. 发送给客户端的消息码前缀：SC_
		3. 接收其他服务器的消息码前缀：S_
		4. 发送给其他服务器的消息码前缀：S2_

 * @author 黄承武

 * @create time 2018.4.23
 */

 package web


 //--接收客户端的消息码

 const MSG_CS_USER_LOGIN int = 1001; //用户登录，请求用户数据

 const MSG_CS_SINGLE int = 1002; //用户请求并上传单关卡信息

 const MSG_CS_RNK int = 1003; //用户请求排行榜数据