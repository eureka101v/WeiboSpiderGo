package main

import (
	"WeiboSpiderGo/config"
	"WeiboSpiderGo/scrapy_rules"
	"WeiboSpiderGo/utils"
	"fmt"
)

var uidLi = utils.GetTargetUidList()

func scrapyInfomation() {
	getInfoC := scrapy_rules.GetDefaultCollector()
	getMoreInfoC := scrapy_rules.GetDefaultCollector()
	scrapy_rules.SetMoreInfoCallback(getMoreInfoC)

	scrapy_rules.SetInfoCallback(getInfoC, getMoreInfoC)

	for _, uid := range uidLi {
		url := fmt.Sprintf("%s/%s/info", scrapy_rules.BaseUrl, uid)
		getInfoC.Visit(url)
	}
	getInfoC.Wait()
	getMoreInfoC.Wait()
}

func scrapyTweet() {
	getTweetsC := scrapy_rules.GetDefaultCollector()
	getContentSubC := scrapy_rules.GetDefaultCollector()
	scrapy_rules.SetFullContentCallback(getContentSubC)
	getCommentSubC := scrapy_rules.GetDefaultCollector()
	scrapy_rules.SetCommentCallback(getCommentSubC)

	scrapy_rules.SetTweetCallback(getTweetsC, getContentSubC, getCommentSubC)

	for _, uid := range uidLi {
		url := scrapy_rules.GetTweetUrl(uid)
		getTweetsC.Visit(url)
	}
	getTweetsC.Wait()
	getContentSubC.Wait()
	getCommentSubC.Wait()
}

func scrapyFollow() {
	getFollowC := scrapy_rules.GetDefaultCollector()
	scrapy_rules.SetFollowCallback(getFollowC)
	//read files
	for _, uid := range uidLi {
		url := scrapy_rules.GetFollowUrl(uid)
		getFollowC.Visit(url)
	}
	getFollowC.Wait()
}

func scrapyFans() {
	getFansC := scrapy_rules.GetDefaultCollector()
	scrapy_rules.SetFansCallback(getFansC)

	for _, uid := range uidLi {
		url := scrapy_rules.GetFansUrl(uid)
		getFansC.Visit(url)
	}
	getFansC.Wait()
}

func main() {
	if config.Conf.GetBool("SCRAPY_TYPE.Info") {
		scrapyInfomation()
	}
	if config.Conf.GetBool("SCRAPY_TYPE.Follow") {
		scrapyFollow()
	}
	//修复去重问题
	if config.Conf.GetBool("SCRAPY_TYPE.Fans") {
		scrapyFans()
	}
	if config.Conf.GetBool("SCRAPY_TYPE.Tweet.Main") {
		scrapyTweet()
	}
}
