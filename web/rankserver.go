/** 
 * 排行榜服务器
 *
 * @author 黄承武
 *
 * @create time 2018.4.24
 */

package web

import(
	"fmt"
	"time"
	"server/sql"
	"server/datastruct"
	"sort"
	"server/config"
)


////////////////////////////////////////////////////////////////////////
//全局变量

//声明、定义一个全局变量。用于存储排行榜排序结果，排序结果不存入数据库，只在服务器缓存中进行
var Rankcache []datastruct.RankResult;


///////////////////////////////////////////////////////////////////////
//排序相关


type stWrapper datastruct.SortWrapper; //创建别名. Go语言中不运行在基本类型变量中进行重定义，需要进行别名转换。
type sortby func(p, q *datastruct.RankResult) bool; //排序函数

func (sw stWrapper) Len() int { // 重写 Len() 方法
    return len(sw.Infos)
}

func (sw stWrapper) Swap(i, j int){ // 重写 Swap() 方法
    sw.Infos[i], sw.Infos[j] = sw.Infos[j], sw.Infos[i]
}

func (sw stWrapper) Less(i, j int) bool { // 重写 Less() 方法
    return sw.By(&sw.Infos[i], &sw.Infos[j])
}

func sortUser(u [] datastruct.RankResult, by sortby){
    sort.Sort(stWrapper{u, by})
}

////////////////////////////////////////////////////////////////////////



//启动排行榜服务
func StartRankServer(){
	pmsg1 := "-> 启动排行榜（子）服务...";
	fmt.Println(pmsg1);
	Pushlog(pmsg1);

	//启动一个线程去执行此操作
	go func(){
		for{
			//执行排行榜操作
			sortRnk();
			
			//获取当前时间
			now := time.Now();

			//将现在的时间加上24小时变成第2天的这个时候，然后根据第2天的当前时间创建一个新的时间为第2天的设定时间
			next := now.Add(time.Hour * 24);
			next = time.Date(next.Year(), next.Month(), next.Day(), config.RANK_REFRESH_TIME, 0, 0, 0, next.Location());

			//将定好的更新时间-当前时间=剩余倒计时时间，启动计时器管道
			t := time.NewTimer(next.Sub(now));

			//发送管道并阻塞子线程，等待管道抛出。
			<-t.C;
		}

	}();
}



//执行排行榜操作
func sortRnk(){
	fmt.Println("-> 更新排行榜.");
	Pushlog("执行排行榜数据的重排操作.");

	//查询数据
	ary := sql.GetAllUserBaseinfo();
	if ary == nil {
		fmt.Println("执行排行榜数据重排时异常! 数据库中暂无任何用户基础信息数据。");
		return;
	}

	// fmt.Println(ary);
	// fmt.Print("\n");

	//进行排序
	sortUser(ary, func(p, q *datastruct.RankResult) bool{
		return p.Maxpoint > q.Maxpoint;
	});

	// fmt.Println(ary);
	// fmt.Print("\n");

	//排序结束后，取出配置中配置的名数
	if len(ary) >= config.RANK_NUMBER {
		Rankcache = ary[:config.RANK_NUMBER];
	}else{
		Rankcache = ary[0:];
	}

	// fmt.Println(Rankcache);
}