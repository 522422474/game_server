/**
 * 日志操作

 * @author 黄承武

 * @create time 2018.4.26
 */

package web


import (
	"strings"
	"fmt"
	"time"
	"os"
	"server/config"
)



/** 日志缓存 */
var logCache []string;

/** 锁 */
var islock bool = false;




/**
	启动服务器日志线程
  */
func StartLogServer(){
	pmsg1 := "-> 启动日志（子）服务...";
	fmt.Println(pmsg1);
	Pushlog(pmsg1);

	//创建本地日志文件目录
	_, err := os.Stat(config.LOG_DIR);
	if err != nil{
		err := os.Mkdir(config.LOG_DIR, os.ModePerm);
		if err != nil {
			fmt.Printf("-> [Error] : 创建日志目录失败![%v]\n", err);
			return;
		}
	}

	//启动线程
	go func() {
		for{
			llen := len(logCache);
			if llen != 0 && !islock{
				fmt.Println("-> 写日志...");
				Pushlog("写日志.");

				//取出N条日志记录写入本地
				var data []string;
				if len(logCache) >= config.LOG_WRITE_NUM {
					data = logCache[0:config.LOG_WRITE_NUM];
					logCache = logCache[config.LOG_WRITE_NUM:];
				}else{
					data = logCache[0:];
					logCache = []string{};
				}

				__log := strings.Join(data, "\r\n");
				// fmt.Println(data, logCache, "\n", __log);

				//写入
				tm := time.Unix(time.Now().Unix(), 0);
				filename := config.LOG_DIR + tm.Format("2006-01-02")+".log";
				f, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0766);
				n, _ := f.Seek(0, os.SEEK_END);
				f.WriteAt([]byte(time.Now().Format("2006-01-02 15:04:05") + "\r\n" + __log + "\r\n\r\n"), n);
				defer f.Close();
			}

			//定时器
			time.Sleep(time.Duration(config.LOG_WRITE_SEQTIME) * time.Minute);
		}
	}();
}




/**
	推送日志内容
	@param msg
  */
func Pushlog(msg string){
	if msg=="" { return; };
	islock = true;

	logCache = append(logCache, msg);

	islock = false;
}