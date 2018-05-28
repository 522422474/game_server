/**
	处理接收客户端的消息

	@author 黄承武

	@create time 2018.4.23

 */


 package web

 import(
	 "net/http"
	 "server/datastruct"
	 "fmt"
	 "server/sql"
	 "encoding/json"
	 "server/utils"
	 "strings"
	 "strconv"
	 "time"
 )

 /** 
	 具体的处理函数 
	 @param w
	 @param obj json数据对象（客户端发送的） 
   */
 func DoHttpHandler(w http.ResponseWriter, obj datastruct.Netmsg){
	// fmt.Println("接收消息：", obj.OpenId, obj.Id, obj.Cmd, obj.Content, obj.Sign);

	
	switch obj.Cmd{

		case MSG_CS_USER_LOGIN: //用户登录消息，数据序列化
		doUserLoginMsg(w, obj);

		case MSG_CS_SINGLE: //单关卡通关后上传，上报
		doUserSingleMsg(w, obj);

		case MSG_CS_RNK: //排行榜数据
		doRnkMsg(w, obj);

		default:
			pmsg := "无效的CMD指令. cmd = " + strconv.Itoa(obj.Cmd);
			fmt.Println(pmsg);
			Pushlog(pmsg);
	}

 }




//--------------------------------------------------------------
//--以下为处理客户端请求的消息



//处理用户登录请求
//@param w
//@param obj
 func doUserLoginMsg(w http.ResponseWriter, obj datastruct.Netmsg){
	/**
		判断登录的用户是否在数据缓存中。
		如果没有：添加到数据缓存中，并给用户返回相关信息（序列化）
		如果有：直接给用户返回用户基础信息（序列化）
	 */
	 val := sql.GetUserBaseinfo(obj.OpenId);
	 var uid int;
	 var maxp int;

	 if val == nil {
		 //没有此用户，添加新用户
		 maxp = 0;
		 uid= utils.BuildID(); //给用户生成一个UID，游戏中的唯一ID
		 obj.Id = uid;
		//  fmt.Println("create",uid, maxp);

		 sql.AddUserBaseinfo(obj);
	 }else{
		 //存在用户，直接给客户端返回数据
		 data := datastruct.UserBaseinfo{};
		 json.Unmarshal(val, &data);
		 uid = data.Id;
		 maxp = data.Maxpoint;
		//  fmt.Println("get", uid, maxp);
		sql.UpdateUserBaseinfo(obj.OpenId, obj.Content, time.Now());
	 }

	 Pushlog("-> 用户登陆：openid = " + obj.OpenId + " / id = " + strconv.Itoa(obj.Id) + " / name = " + obj.Content + " / time = " + time.Now().Format("2006-01-02 15:04:05"));


	 msg := datastruct.LoginResult{};
	 msg.Maxpoint = maxp;
	 msg.Uid = uid;
	 msgdata, _ := json.Marshal(msg);
	 fmt.Fprintf(w, string(msgdata));
 }



 //处理接收客户端单关卡上报消息
 //@param w
 //@param obj
 func doUserSingleMsg(w http.ResponseWriter, obj datastruct.Netmsg){
	// fmt.Println(obj.Content);
	ary := strings.Split(obj.Content, ",");
	point := ary[0];
	tim := ary[1];
	
	//获取用户基础信息
	val := sql.GetUserBaseinfo(obj.OpenId);
	if val == nil {
		pmsg1 := "处理用户上传单关卡数据信息发生异常，此用户不在数据缓存中. OpenId=" + obj.OpenId;
		fmt.Println(pmsg1);
		Pushlog(pmsg1);
		return;
	}
	
	data := datastruct.UserBaseinfo{};
	json.Unmarshal(val, &data);

	sendpoint, _ := strconv.Atoi(point);
	sendtime, _ := strconv.Atoi(tim);
	
	isUpdateBaseinfo := false;
	tm := time.Now();
	t1 := strings.Split(tm.Format("2006-01-02 15:04:05"), " ")[0];
	t2 := strings.Split(data.Logintime.Format("2006-01-02 15:04:05"), " ")[0];
	if t1 != t2 {
		isUpdateBaseinfo=true;
		data.Logintime = tm; //这里的作用是过凌晨12点时，用户一直保持连接，并没有重新登陆，那么通过上报关卡数据时更新登陆时间.进一步精确DAU日登陆数据
		if data.NewPlayer {
			data.NewPlayer = false;
		}
	}
	// fmt.Println(sendpoint, sendtime, tm);

	if data.Maxpoint < sendpoint { //当前用户最大关卡数小于用户上报关卡数
		//更新基础信息
		data.Maxpoint = sendpoint;
		//创建一条新的关卡记录
		single := datastruct.UserSingleCardinfo{};
		single.OpenId = data.OpenId;
		single.Id = data.Id;
		single.Point = sendpoint;
		single.Time = sendtime;
		single.Reporttime = tm;
		sql.UpdateUserMaxpointAndInsertnew(data, single);

	}else{ //当前用户最大关卡数大于等于上报关卡数
		sql.UpdateUserSinglecard(obj.Id, sendpoint, sendtime, isUpdateBaseinfo, data); 
	}

	//给客户端发送一个服务器收到此消息，但没有任何返回值。
	ret := datastruct.NullMessage{};
	ret.Code = 0;
	retMsg, _ := json.Marshal(ret);
	fmt.Fprintln(w, string(retMsg));
 }



 //处理客户端请求排行榜信息消息
 //@param w
 //@param obj
 func doRnkMsg(w http.ResponseWriter, obj datastruct.Netmsg){
	// fmt.Println("do rank message.");

	//创建数据对象
	rs := datastruct.RankResultToClient{};

	if Rankcache == nil {
		fmt.Println("服务器排行榜暂无数据，等待刷新时间进行重排。");
		rs.Source = []datastruct.RankResult{}

	}else{
		rs.Source = Rankcache;
	}

	data, _ := json.Marshal(rs);
	fmt.Fprintln(w, string(data));
 }