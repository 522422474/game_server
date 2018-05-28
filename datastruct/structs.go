/**
   数据结构体管理
   管理、封装服务器中定义的结构体对象
 
   @author 黄承武
   
   @create time 2018.4.21
 */

package datastruct
 
import (
	"time"
)


//接收客户端消息数据结构（客户端、服务器必须按照此结构，否则会报错或无法解析正确数据）
type Netmsg struct
{
	OpenId string //平台账号
	Id int //用户ID
	Cmd int //消息处理指令（详见msgcode.go）
	Content string //消息内容（字符串格式，每个数据用“，”分割）
	Sign string //消息验证码（通过客户端与服务器制定的统一算法，来判断是否是合法消息. ）
}



//用户基础信息数据结构（redis数据库中用户基础信息储存结构）
 type UserBaseinfo struct
 {
	 OpenId string //平台账号
	 Id int //唯一索引
	 Name string //用户昵称
	 Maxpoint int //最大关卡数
	 Createtime time.Time //用户创建时间
	 Logintime time.Time //最近登陆时间
	 NewPlayer bool //是否是新用户（true：是；false：不是）
 }



//用户单个关卡数据（redis数据库中用户单关卡的信息储存结构）
 type UserSingleCardinfo struct
 {
	 OpenId string //平台账号
	 Id int //用户唯一索引
	 Point int //关卡id
	 Time int //关卡用时（单位：秒，最大记录值为999。超过此值为999.）
	 Reporttime time.Time //关卡上报提交时间
 }



//用于用户登陆时返回的结果对象结构
type LoginResult struct
{
	Uid int //用户在游戏中唯一ID
	Maxpoint int //用户最大通关数
}



//用于排行榜排序数据对象的结构
type SortWrapper struct
{
	Infos [] RankResult //用户基础信息组
	By func(p, q * RankResult) bool //排序函数
}



//返回给客户端的排行榜单条数据对象结构
type RankResult struct
{
	Name string //昵称
	Maxpoint int //用户最大通关数
}



//发送给客户端的排行榜数据结构
type RankResultToClient struct
{
	Source []RankResult //数据组
}



//空消息数据结构体
type NullMessage struct
{
	Code int //空消息码
}