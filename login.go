package main

import (
	"WeiboSpiderGo/config"
	"WeiboSpiderGo/mdb"
	"bufio"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"gopkg.in/mgo.v2/bson"
	"os"
	"strings"
)

var db_name = config.Conf.GetString("DB_NAME")

func Displayed(by, elementName string) func(selenium.WebDriver) (bool, error) {
	return func(wd selenium.WebDriver) (bool, error) {
		el, err := wd.FindElement(by, elementName)
		if err != nil {
			return false, nil
		}
		enabled, err := el.IsDisplayed()
		if err != nil {
			return false, nil
		}

		if !enabled {
			return false, nil
		}

		return true, nil
	}
}

func getCookieStr(username_text string, password_text string) string {
	// Start a Selenium WebDriver server instance (if one is not already
	// running).
	//username_text := "1222"
	//password_text := "23121"
	const (
		// These paths will be different on your system.
		driverPath = "/usr/local/bin/chromedriver"
		port       = 9005
	)
	opts := []selenium.ServiceOption{}

	service, err := selenium.NewChromeDriverService(driverPath, port, opts...)
	if nil != err {
		fmt.Println("start a chromedriver service falid", err.Error())
		return ""
	}
	//注意这里，server关闭之后，chrome窗口也会关闭
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			//"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模拟user-agent，防反爬
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		fmt.Println("connect to the webDriver faild", err.Error())
		return ""
	}
	defer wd.Quit()
	err = wd.Get("https://passport.weibo.cn/signin/login?entry=mweibo&r=https://weibo.cn/")
	if err != nil {
		fmt.Println("get page faild", err.Error())
		return ""
	}
	wd.Wait(Displayed(selenium.ByCSSSelector, "#loginName"))
	wd.Wait(Displayed(selenium.ByCSSSelector, "#loginPassword"))
	wd.Wait(Displayed(selenium.ByCSSSelector, "#loginAction"))
	username, err := wd.FindElement(selenium.ByCSSSelector, "#loginName")
	if err != nil {
		fmt.Println("get username faild", err.Error())
		return ""
	}
	password, err := wd.FindElement(selenium.ByCSSSelector, "#loginPassword")
	if err != nil {
		fmt.Println("get username faild", err.Error())
		return ""
	}
	submit, err := wd.FindElement(selenium.ByCSSSelector, "#loginAction")
	if err != nil {
		fmt.Println("get username faild", err.Error())
		return ""
	}
	username.SendKeys(username_text)
	password.SendKeys(password_text)
	submit.Click()
	wd.Wait(func(wdtemp selenium.WebDriver) (b bool, e error) {
		tit, err := wdtemp.Title()
		if err != nil {
			return false, nil
		}
		if tit != "我的首页" {
			return false, nil
		}
		return true, nil
	})
	mcookie, err := wd.GetCookies()
	var cookie_arr []string
	for _, c := range mcookie {
		cookie_arr = append(cookie_arr, c.Name+"="+c.Value)
	}
	cookie_str := strings.Join(cookie_arr, ";")
	return cookie_str
}

func saveToMgo(id_ string, password string, cookie_str string) {
	err := mdb.Upsert(db_name, "account", bson.M{"_id": id_}, bson.M{"$set": bson.M{"password": password, "cookie": cookie_str, "status": "success"}})
	if err != nil{
		panic(err)
	}
	fmt.Println("login success")
}

func main() {
	file, err := os.Open(config.Conf.GetString("ACCOUNT_FILE"))
	if err != nil {
		fmt.Println("account.txt is not found")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		text := strings.Split(lineText, "----")
		fmt.Println("start login username:",text[0])
		cookiestr := getCookieStr(text[0],text[1])
		saveToMgo(text[0],text[1],cookiestr)
	}
}
