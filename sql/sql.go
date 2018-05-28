/**
	数据库redis相关操作
	注:	redis 0库 ：负责存储用户的基础信息
		redis 1库 ：负责存储每个用户的关卡详细记录

	@author 黄承武

	@create time 2018.4.23

 */
package sql


import(
	"github.com/garyburd/redigo/redis"
	"server/datastruct"
	"fmt"
	"encoding/json"
	"server/config"
	"time"
	"strings"
)

//  var connPools chan redis.Conn;
 
//  func pushConn(conn redis.Conn) {
// 	 // 基于函数和接口间互不信任原则，这里再判断一次
// 	 if connPools == nil {
// 		 connPools = make(chan redis.Conn, config.SQL_CONN_MAX_LENGTH);
// 	 }
// 	 if len(connPools) >= config.SQL_CONN_MAX_LENGTH {
// 		 conn.Close();
// 		 return;
// 	 }
// 	 connPools <- conn;
//  }
 
//  func getConn(network string, address string) redis.Conn {
// 	 // 缓冲机制，相当于消息队列
// 	 if len(connPools) == 0 {
// 		 // 如果长度为0，就定义一个redis.Conn类型长度为config.SQL_CONN_MAX_LENGTH的channel
// 		 connPools = make(chan redis.Conn, config.SQL_CONN_MAX_LENGTH);
// 		 //创建子线程
// 		 go func() {
// 			 for i := 0; i < config.SQL_CONN_MAX_LENGTH/2; i++ {
// 				 c, err := redis.Dial(network, address); //链接数据库（tcp协议）
// 				 if err != nil {
// 					 panic(err);
// 				 }
// 				 pushConn(c);
// 			 }
// 		 } ()
// 	 }
// 	 return <-connPools
//  }


 /**
  *	数据库连接池
  */
 var pool *redis.Pool
 func newPool() *redis.Pool {
	 return &redis.Pool {
		 MaxIdle:     config.SQL_MAXIDLE,
		 MaxActive:   config.SQL_MAXACTIVE,
		 IdleTimeout: 300 * time.Second,
		 Dial: func() (redis.Conn, error) {
			 c, err := redis.Dial("tcp", config.SQL_URL_PORT);
			 if err != nil {
				 return nil, err
			 }
			 if _, err := c.Do("AUTH", config.SQL_PASSWORD); err != nil {
				 c.Close()
				 return nil, err
			 }
			 return c, err
		 },
		 TestOnBorrow: func(c redis.Conn, t time.Time) error {
			 if time.Since(t) < time.Minute {
				 return nil
			 }
			 _, err := c.Do("PING");
			 return err
		 },
	 }
 }


 //------------------------------------------------

 
 /**
	 向数据库中添加用户
	 @param data 消息数据对象
   */
func AddUserBaseinfo(data datastruct.Netmsg) (bool) {
	
	//创建数据
	table := datastruct.UserBaseinfo{};
	table.OpenId = data.OpenId;
	table.Id = data.Id;
	table.Name = data.Content;
	table.Maxpoint = 0;
	table.Createtime = time.Now();
	table.Logintime = time.Now();
	table.NewPlayer = true;
	userdata, _ := json.Marshal(table);

	//链接
	c := newPool().Get();
	defer c.Close();
   
	//选择redis库
	c.Do("select", 0);

	//添加用户
	_, err := c.Do("rpush", data.OpenId, userdata);

	var result bool;
	if err != nil{
		result = false;
	}else{
		result = true;
	}

	return result;
}



/**
	更新用户基础信息
	@param openid 平台账号
	@param newname 新的昵称
	@param t 最新登陆时间
*/
func UpdateUserBaseinfo(openid string, newname string, t time.Time){
	val := GetUserBaseinfo(openid);
	if(val == nil){ return;};
	data := datastruct.UserBaseinfo{};
	json.Unmarshal(val, &data);

	//---------------------------
	//--需要更新的信息如下
	
	data.Name = newname;
	data.Logintime = t;
	
	crttime := data.Createtime.Format("2006-01-02 15:04:05");
	timestr := t.Format("2006-01-02 15:04:05");
	if strings.Split(crttime, " ")[0] == strings.Split(timestr, " ")[0] {
		//最近登陆时间与创建用户的时间是一致的，认为还是新用户
		data.NewPlayer = true;
	}else{
		//老用户
		data.NewPlayer = false;
	}
	
	//--END
	//---------------------------

	c := newPool().Get();
	defer c.Close();

	c.Do("select", 0);

	newval, _ := json.Marshal(data);
	c.Do("lset", openid, 0, newval);
}



/**
	获取用户基础信息
	@param openid 平台唯一账号
 */
 func GetUserBaseinfo(openid string)( []byte ){    

	c := newPool().Get();
	defer c.Close();

	c.Do("select", 0);
	values, err := c.Do("lindex", openid, "0");

	if err != nil  {
		fmt.Println(err);
	}

	if(values == nil){
		return nil;
	}
	return values.([]byte);
}



/**
	获取用户最大关卡数
	@param openid 平台账号
*/
func GetUserMaxpoint(openid string) (int){
	val := GetUserBaseinfo(openid);
	if(val == nil){
		return 0;
	}
	data := datastruct.UserBaseinfo{};
	json.Unmarshal(val, &data);
	return data.Maxpoint;
}



/**
 * 获取所有用户的基础信息（用于世界榜排名操作）
 */
 func GetAllUserBaseinfo() ( []datastruct.RankResult ) {
	c := newPool().Get();
	defer c.Close();

	//选择基础信息库
	c.Do("select","0");

	//获取所有用户的的key值
	data, err := redis.Values(c.Do("keys","*"));
	if err != nil{
		fmt.Println(err);
	}

	if data == nil {
		return nil;
	}

	
	//利用key值获取相对应的值，并储存返回
	ary := []datastruct.RankResult{};

	for _, v := range data{

		baseinfo, _ := c.Do("lindex", string(v.([]byte)), 0);

		bi := datastruct.UserBaseinfo{};
		json.Unmarshal(baseinfo.([]byte), &bi);

		result := datastruct.RankResult{};
		result.Maxpoint = bi.Maxpoint;
		result.Name = bi.Name;

		ary = append(ary, result);
	}


	return ary;
}



//创建用户唯一ID
//@return 创建的唯一ID 类型int64
func GetID()(int64){
	c := newPool().Get();
	defer c.Close();

	//选择库
	_, err := c.Do("select", 0);
	if err!=nil{
		return -1;
	}

	val, err := c.Do("dbsize");
	if err!=nil {
		return -1;
	}

	return val.(int64);
}



/**
	更新用户最大关卡数，并插入新的一条关卡信息
	@param obj 用户基础信息
	@param newinfo 用户单关信息
 */
func UpdateUserMaxpointAndInsertnew(obj datastruct.UserBaseinfo, newinfo datastruct.UserSingleCardinfo){
	c := newPool().Get();
	defer c.Close();

	//更新用户基础信息
	c.Do("select", 0);
	val, _ := json.Marshal(obj);
	c.Do("lset", obj.OpenId, 0, val);

	//插入一条新的信息
	newval, _ := json.Marshal(newinfo);
	c.Do("select", 1);
	c.Do("rpush", obj.Id, newval);
}



//更新用户某一关卡信息
//@param id 用户ID
//@param point 关卡索引（从0开始）
//@param time 当前关卡用时（单位：秒）
//@param isupdateBaseinfo 是否更新用户基础信息
//@param obj 用户基础信息数据
func UpdateUserSinglecard(id int , point int, time int, isupdateBaseinfo bool, obj datastruct.UserBaseinfo){
	c := newPool().Get();
	defer c.Close();

	c.Do("select", 1);

	cardIndex := point -1;

	values, err:= c.Do("lindex", id, cardIndex);
	if values == nil {
		fmt.Println("更新用户某一关卡用时操作发生异常！未找到当前用户当前关卡数据。id=", id, " / point=", point);
		return; 
	}

	if err != nil {
		fmt.Println(err);
	}

	data := datastruct.UserSingleCardinfo{};
	json.Unmarshal(values.([]byte), &data);

	if time < data.Time {
		//上报的用时必须小于现有用时，才更新
		data.Time = time;
		info, _ := json.Marshal(data);
		c.Do("lset", id, cardIndex, info);
	}

	if isupdateBaseinfo {
		c.Do("select", 0);
		bi, _ := json.Marshal(obj);
		c.Do("lset", obj.OpenId, 0, bi);
	}
}