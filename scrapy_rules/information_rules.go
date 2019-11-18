package scrapy_rules

import (
	"WeiboSpiderGo/mdb"
	"WeiboSpiderGo/utils"
	"fmt"
	"github.com/gocolly/colly"
	"gopkg.in/mgo.v2"
	"strconv"
	"strings"
	"time"
)

func SetInfoCallback(getInfoC, getMoreInfoC *colly.Collector) {
	getInfoC.OnResponse(func(r *colly.Response) {
		content := string(r.Body)
		info := mdb.Information{}
		info.CrawlTime = int32(time.Now().Unix())
		info.Id_ = utils.ReParse(`(\d+)/info`, r.Request.URL.String())
		nickName := utils.ReParse(`昵称?[：:]?(.*?)<br/>`, content)
		authentication := utils.ReParse(`认证?[：:]?(.*?)<br/>`, content)
		gender := utils.ReParse(`性别?[：:]?(.*?)<br/>`, content)
		place := utils.ReParse(`地区?[：:]?(.*?)<br/>`, content)
		briefIntroduction := utils.ReParse(`简介?[：:]?(.*?)<br/>`, content)
		birthday := utils.ReParse(`生日?[：:]?(.*?)<br/>`, content)
		sexOrientation := utils.ReParse(`性取向?[：:]?(.*?)<br/>`, content)
		sentiment := utils.ReParse(`感情状况?[：:]?(.*?)<br/>`, content)
		vipLevel := utils.ReParse(`会员等级?[：:]?(.*?)&nbsp;<a`, content)
		labels := utils.ReParseMayLi(`stag=1">(.*?)</a>`, content) //标签
		info.Nickname = nickName
		info.Gender = gender
		placeli := strings.Split(place, " ")
		info.Province = placeli[0]
		if len(placeli) > 1 {
			info.City = placeli[1]
		}
		info.BriefIntroduction = briefIntroduction
		info.Birthday = birthday
		if sexOrientation == gender {
			info.SexOrientation = "同性恋"
		} else {
			info.SexOrientation = "异性恋"
		}
		info.Sentiment = sentiment
		info.VipLevel = vipLevel
		info.Authentication = authentication
		info.Labels = ""
		for i, labelItem := range labels {
			if i != 0 {
				info.Labels += ","
			}
			info.Labels += labelItem[1]
		}
		r.Ctx.Put("info", info)
		getMoreInfoC.Request("GET","https://weibo.cn/u/" + info.Id_,nil,r.Ctx,nil)
	})
}

func SetMoreInfoCallback(getMoreInfoC *colly.Collector){
	getMoreInfoC.OnResponse(func(r *colly.Response) {
		content := string(r.Body)
		info := r.Ctx.GetAny("info").(mdb.Information)
		tweetsNum := utils.ReParse(`微博\[(\d+)\]`, content)
		followsNum := utils.ReParse(`关注\[(\d+)\]`, content)
		fansNum := utils.ReParse(`粉丝\[(\d+)\]`, content)
		if tweetsNum != "" {
			temp, _ := strconv.Atoi(tweetsNum)
			info.TweetsNum = int32(temp)
		}
		if followsNum != "" {
			temp, _ := strconv.Atoi(followsNum)
			info.FollowsNum = int32(temp)
		}
		if fansNum != "" {
			temp, _ := strconv.Atoi(fansNum)
			info.FansNum = int32(temp)
		}
		err := mdb.Insert(dbName, "Information", info)
		if mgo.IsDup(err) {
			//有重复数据
			fmt.Println("already scrapy")
		}
	})
}
