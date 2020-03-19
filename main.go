package main

import (
	_ "classsOne/routers"
	"github.com/astaxie/beego"
	_"classsOne/models"
	"strconv"
	//"github.com/astaxie/beego/plugins/cors"
)


func main() {


	beego.AddFuncMap("ShowPrePage",HandlePrePage)//视图函数的映射关系一定要放在run前
	beego.AddFuncMap("ShowNextPage",HandleNextPage)

	beego.Run()

}

//这个处理函数可以放在任何位置
func HandlePrePage(data string)(string)  {
	dataTemp, _ := strconv.Atoi(data)
	pageIndex := dataTemp - 1
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1
}

//这个处理函数可以放在任何位置
func HandleNextPage(data string)(string)  {
	dataTemp, _ := strconv.Atoi(data)
	pageIndex := dataTemp + 1
	pageIndex1 := strconv.Itoa(pageIndex)
	return pageIndex1

}