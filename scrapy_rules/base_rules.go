package scrapy_rules

import (
	"WeiboSpiderGo/config"
	"WeiboSpiderGo/mdb"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"net/http"
	"time"
)

var BaseUrl = "https://weibo.cn"
var cookieStrLi []mdb.Account
var dbName = config.Conf.GetString("DB_NAME")

// return a collector
func GetDefaultCollector() *colly.Collector {
	//set async and dont forget set c.wait()
	if config.Conf.GetBool("DEBUG_MODE") {
	}
	debugger := &debug.LogDebugger{}

	//file,err := os.Create(utils.ExecPath+"/debug.log")
	//if err!=nil{
	//	panic(err)
	//}
	//debugger.Output = file

	var c = colly.NewCollector(
		colly.Async(true),
		colly.Debugger(debugger),
	)
	//disable http KeepAlives its could cause OOM in long time work
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})
	mdb.FindAll(dbName, "account", bson.M{}, bson.M{}, &cookieStrLi)
	setDefaultCallback(c)
	extensions.RandomUserAgent(c)
	return c
}

// set default call,cookie and error handling
func setDefaultCallback(c *colly.Collector) {
	// set random cookie
	c.OnRequest(func(r *colly.Request) {
		n := rand.Intn(len(cookieStrLi))
		r.Headers.Set("Cookie", cookieStrLi[n].Cookie)
		r.Ctx.Put("_id", cookieStrLi[n].Id_)
	})

	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.

	// delay 3 to 5 second
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2, Delay: 8 * time.Second, RandomDelay: 2 * time.Second})

	// deal with error statusCode
	c.OnError(func(r *colly.Response, e error) {
		if r.StatusCode == 302 || r.StatusCode == 403 {
			mdb.Update(dbName, "account", bson.M{"_id": r.Ctx.Get("_id")}, bson.M{"$set": bson.M{"status": "error"}})
		} else if r.StatusCode == 418 {
			fmt.Println("please wait a second")
		}
	})
}
