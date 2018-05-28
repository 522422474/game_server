/**
   服务器配置文件
   配置http服务的地址、端口号，数据库相关信息以及缓存服务的相关信息

   @author 黄承武

   @create time 2018.4.23

  */
package config



//http服务地址
const HTTP_URL string ="192.168.2.254";

//http服务地址的端口号
const HTTP_PORT int = 1201;



//数据库地址、端口号
 const SQL_URL_PORT string = "192.168.2.254";

//数据库登陆密码
const SQL_PASSWORD string = "@doinggame2018!";

//数据库链接管道最大缓冲值 --（暂时废弃）
const SQL_CONN_MAX_LENGTH int = 20;

//最大的激活连接数，表示同时最多有N个连接 ，为0事表示没有限制
const SQL_MAXACTIVE int = 20;

//最大的空闲连接数，表示即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
const SQL_MAXIDLE int = 10;



//排行榜刷新时间（现规则制定为每天凌晨5点）
const RANK_REFRESH_TIME int = 5;

//排行榜取前N名
const RANK_NUMBER int = 100;



//存储日志的路径
const LOG_DIR string = "log/";

//每次写入日志数量
const LOG_WRITE_NUM int = 10000;

//写日志的时间间隔(单位：分钟)
const LOG_WRITE_SEQTIME int = 1;
