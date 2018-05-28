/**
 * 服务器主文件

 * @auhtor 黄承武

 * @create time 2018.4.23
 */


package main

import (
	"server/web"
	"runtime"
)


/**主函数入口*/
func main(){

	//启动Go语言最佳性能，创建的线程数量使用CPU的数量.这样的话Go创建的协程
	//程序的线程不会超过CPU，在线程切换时CPU不会有太多的消耗。
	runtime.GOMAXPROCS(runtime.NumCPU());

	//启动http服务器（主线程）
	web.StartHttpServer(true);

}