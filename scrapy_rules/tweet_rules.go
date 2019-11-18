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

func SetTweetCallback(getTweetsC, getContentSubC, getCommentSubC *colly.Collector) {
	getTweetsC.OnResponse(func(r *colly.Response) {
		content := string(r.Body)
		uid := utils.ReParse(`(\d+)/profile`, r.Request.URL.String())
		if strings.Contains(r.Request.URL.String(), "page=1") {
			allPage := utils.ReParse(`/>&nbsp;1/(\d+)页</div>`, content)
			pageNum, _ := strconv.Atoi(allPage)
			for i := 2; i < (pageNum + 1); i++ {
				link := fmt.Sprintf("%s/%s/profile?page=%d",BaseUrl,uid,i)
				getTweetsC.Visit(link)
			}
		}
	})
	getTweetsC.OnXML(`//div[@class="c" and @id]`, func(element *colly.XMLElement) {
		tweet := mdb.Tweets{}
		tweet.CrawlTime = int32(time.Now().Unix())
		tweetRepostUrl := element.ChildAttr(`.//a[contains(text(),"转发[")]`, "href")
		tweetItemId := utils.ReParse(`/repost/(.*?)\?`, tweetRepostUrl)
		tweet.UserId = utils.ReParse(`uid=(\d+)`, tweetRepostUrl)
		tweet.WeiboUrl = fmt.Sprintf("https://weibo.com/%s/%s", tweet.UserId, tweetItemId)
		tweet.Id_ = fmt.Sprintf("%s_%s", tweet.UserId, tweetItemId)
		createTimeInfo := element.ChildText(`.//span[@class="ct"]`)
		if strings.Contains(createTimeInfo, "来自") {
			timeStr := strings.Split(createTimeInfo, "来自")[0]
			timeStr = strings.TrimSpace(timeStr)
			tweet.CreatedAt = utils.ConvTime(timeStr)
			tweet.Tool = strings.Split(createTimeInfo, "来自")[1]
		} else {
			timeStr := strings.TrimSpace(createTimeInfo)
			tweet.CreatedAt = utils.ConvTime(timeStr)
		}

		likeNumText := element.ChildText(`.//a[contains(text(),"赞[")]`)
		likeNum, _ := strconv.Atoi(utils.ReParse(`\d+`, likeNumText))
		tweet.LikeNum = int32(likeNum)

		repostNumText := element.ChildText(`.//a[contains(text(),"转发[")]`)
		repostNum, _ := strconv.Atoi(utils.ReParse(`\d+`, repostNumText))
		tweet.RepostNum = int32(repostNum)

		commentNumText := element.ChildText(`.//a[contains(text(),"评论[") and not(contains(text(),"原文"))]`)
		commentNum, _ := strconv.Atoi(utils.ReParse(`\d+`, commentNumText))
		tweet.CommentNum = int32(commentNum)

		tweet.ImageUrl = element.ChildAttr(`.//img[@alt="图片"]`, "src")
		tweet.VideoUrl = element.ChildAttr(`.//a[contains(@href,"https://m.weibo.cn/s/video/show?object_id=")]`, "href")

		mapNode := element.ChildAttr(`.//a[contains(text(),"显示地图")]`, "href")
		if mapNode != "" {
			tweet.LocationMapInfo = utils.ReParse(`xy=(.*?)&`, mapNode)
		}

		tweet.OriginWeibo = element.ChildAttr(`.//a[contains(text(),"原文评论[")]`, "href")

		allContentLink := element.ChildAttr(`.//a[text()="全文" and contains(@href,"ckAll=1")]`, "href")
		if allContentLink == "" {
			//没有全文按钮
			content := element.Text
			if pos := strings.LastIndex(content, "转发理由:"); pos != -1 {
				content = content[pos+len("转发理由:"):]
			}
			content = content[0:strings.LastIndex(content, "赞")]
			if pos := strings.LastIndex(content, "[组图共"); pos != -1 {
				content = content[0:pos]
			}
			if pos := strings.LastIndex(content, "原图"); pos != -1 {
				l := len(content)
				if l >= pos+6 {
					content = content[0:pos]
				}
			}
			tweet.Content = strings.TrimSpace(content)
			err := mdb.Insert(dbName, "Tweets", tweet)
			if mgo.IsDup(err) {
				//有重复数据
				fmt.Println("already scrapy")
			}
		} else {
			element.Response.Ctx.Put("tweet", tweet)
			contentSubLink := fmt.Sprintf("%s%s",BaseUrl,allContentLink)
			getContentSubC.Request("GET", contentSubLink, nil, element.Response.Ctx, nil)
		}

		commentLink := fmt.Sprintf("%s/comment/%s?page=1",BaseUrl,strings.Split(tweet.Id_, "_")[1])
		element.Response.Ctx.Put("weibo_url", tweet.WeiboUrl)
		getCommentSubC.Request("GET", commentLink, nil, element.Response.Ctx, nil)
	})
}

func SetFullContentCallback(getContentSubC *colly.Collector) {
	getContentSubC.OnXML(`//*[@id="M_"]/div[1]`, func(element *colly.XMLElement) {
		//var tweet mdb.Tweets
		tweetInt := element.Response.Ctx.GetAny("tweet")
		tweet := tweetInt.(mdb.Tweets)
		content := element.Text
		if pos := strings.LastIndex(content, "转发理由:"); pos != -1 {
			content = content[pos+len("转发理由:"):]
		}
		if pos := strings.LastIndex(content, "[组图共"); pos != -1 {
			content = content[0:pos]
		}
		if pos := strings.LastIndex(content, "原图"); pos != -1 {
			l := len(content)
			if l >= pos+6 {
				content = content[0:pos]
			}
		}
		tweet.Content = strings.TrimSpace(content)
		err := mdb.Insert(dbName, "Tweets", tweet)
		if mgo.IsDup(err) {
			//有重复数据
			fmt.Println("already scrapy")
		}
	})
}

func SetCommentCallback(getCommentSubC *colly.Collector) {
	getCommentSubC.OnResponse(func(r *colly.Response) {
		content := string(r.Body)
		if strings.Contains(r.Request.URL.String(), "page=1") {
			allPage := utils.ReParse(`/>&nbsp;1/(\d+)页</div>`, content)
			pageNum, _ := strconv.Atoi(allPage)
			for i := 2; i < (pageNum + 1); i++ {
				pageUrl := strings.Replace(r.Request.URL.String(), "page=1", "page="+strconv.Itoa(i), -1)
				getCommentSubC.Visit(pageUrl)
			}
		}
	})
	getCommentSubC.OnXML(`//div[@class="c" and contains(@id,"C_")]`, func(element *colly.XMLElement) {
		commentUserUrl := element.ChildAttr(`.//a[contains(@href,"/u/")]`, "href")
		if commentUserUrl == "" {
			return
		}
		comment := mdb.Comment{}
		comment.CrawlTime = int32(time.Now().Unix())
		comment.WeiboUrl = element.Response.Ctx.Get("weibo_url")
		comment.CommentUserId = utils.ReParse(`/u/(\d+)`, commentUserUrl)
		comment.Id_ = element.Attr("id")
		createdAtInfo := element.ChildText(`.//span[@class="ct"]`)
		likeNumText := element.ChildText(`.//a[contains(text(),"赞[")]`)
		likeNum, _ := strconv.Atoi(utils.ReParse(`\d+`, likeNumText))
		comment.LikeNum = int32(likeNum)
		comment.CreatedAt = utils.ConvTime(strings.Split(createdAtInfo, "\u0000")[0])
		content := element.Text
		content = content[0:strings.LastIndex(content, "赞")]
		if pos := strings.LastIndex(content, "举报"); pos != -1 {
			content = content[0:pos]
		}
		comment.Content = strings.TrimSpace(content)
		err := mdb.Insert(dbName, "Comments", comment)
		if mgo.IsDup(err) {
			//有重复数据
			fmt.Println("already scrapy")
		}
	})
}

func GetTweetUrl(uid string)string{
	return fmt.Sprintf("%s/%s/profile?page=1",BaseUrl,uid)
}
