/**
   工具类

   @author 黄承武

   @create time 2018.4.23

 */
package utils

import(
	"crypto/md5"
	"strconv"
	"encoding/hex"
	// "math/rand"
	// "time"
	"server/sql"
)




/**
	检测客户端发送的sign值是否与服务器计算的相等
	@param c_sign
	@param uid
	@param cmd
	@param cont
	@return true：相等；false：不相等（非法）
  */
func CheckMessageSign(c_sign string, openid string,  uid int, cmd int, cont string) (bool){
	
	/* 
		 计算格式如下：
		 MD5（ 应用程序ID（这里自定义dl） + 用户ID + 用户发送的操作指令 + 用户发送的消息。分隔符用“@” ）
	 */
	key := ("dl" + "@" + openid + "@" + strconv.Itoa(uid) + "@" + strconv.Itoa(cmd) + "@" + cont);
	// fmt.Println(key);

	ctx := md5.New();
	ctx.Write([]byte(key));
	str := ctx.Sum(nil);
	md5str := hex.EncodeToString(str);
	// fmt.Println(md5str);
	
	if md5str == c_sign {
		return true;
	}

	return false;
}



//创建唯一ID
//@return 唯一ID
func BuildID()(int){
	// rad := rand.New(rand.NewSource(time.Now().Unix()));
	// return rad.Intn(100000000);
	return int(sql.GetID()+1);
}