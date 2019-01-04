package content

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
)

type Cont struct {
	Title   string
	Content string
	Url     string
}

var baseUrl = "https://m.biqiuge.com/"

// GetLastTitleList 获取最近更新列表
func GetLastTitleList(url string, num int) []Cont {
	var ret []Cont
	doc, _ := get(url)
	doc.Find(".books .listpage select option").Each(func(i int, s *goquery.Selection) {
		values, _ := s.Attr("value")
		c := Cont{Url: values}
		ret = append(ret, c)
	})

	//取最后一页上的所有章节名
	if len(ret) < 1 {
		log.Println("没有取到列表")
		return ret
	}
	lastC := ret[len(ret)-1]
	titleListDoc, _ := get(baseUrl + lastC.Url)
	titleLists := getTitleList(titleListDoc)

	//取最后一页的内容太少，再取前一页，防止一次更新太多时，会缺少
	if len(titleLists) < num {
		secondC := ret[len(ret)-2]
		c, _ := get(baseUrl + secondC.Url)
		titleLists = append(getTitleList(c), titleLists...)
	}

	log.Printf("title-list\n:%+v\n", titleLists)
	return titleLists
}

func GetCont(u string) string {
	u = baseUrl + u
	var ret string
	urlObj, _ := url.Parse(u)
	pathArr := strings.Split(urlObj.Path, ".")
	for i := 1; i < 10; i++ {
		urlTemp := urlObj.Scheme + "://" + urlObj.Host + pathArr[0]
		if i > 1 {
			urlTemp = urlTemp + "_" + strconv.Itoa(i) + ".html"
		} else {
			urlTemp = urlTemp + ".html"
		}
		cont, err := get(urlTemp)
		if err != nil {
			log.Print("%v\n", err)
			return ret
		}
		pageReg := regexp.MustCompile(`.*\(.(\d)/(\d).\).*`)
		params := pageReg.FindStringSubmatch(cont.Text())
		//没有这个分页
		if params[1] > params[2] {
			return ret
		}
		c, _ := cont.Find("#chaptercontent").Html()
		ret = ret + clearCont(c)
	}
	//	fmt.Printf("%v", ret)
	return ret
}

func get(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", url)
	req.Header.Set("User-Agent", " Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1")
	res, err := client.Do(req)
	if err != nil {
		log.Println("get 请求失败")
		log.Printf("%v", err)
		return nil, err
	}
	defer res.Body.Close()

	utfBody, err := iconv.NewReader(res.Body, "gbk", "utf-8")
	if err != nil {
		log.Println("转编码失败")
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		log.Println("页面返回数据为空")
		return nil, err
	}
	return doc, nil
}

func getTitleList(titleListDoc *goquery.Document) []Cont {
	var titleLists []Cont
	titleListDoc.Find(".books .book_last").Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			s.Find("dl dd a").Each(func(ii int, ss *goquery.Selection) {
				urlStr, eUrlStr := ss.Attr("href")
				if eUrlStr {
					t := Cont{Url: urlStr, Title: ss.Text()}
					titleLists = append(titleLists, t)
				}
			})
		}
	})
	return titleLists
}

//清理无用段落
func clearCont(c string) string {

	clearReg := regexp.MustCompile(`<p class="readinline">.*</p>`)
	c = clearReg.ReplaceAllString(c, "")

	clearReg = regexp.MustCompile(`第\d{1,8}.*<br/><br/>`)
	c = clearReg.ReplaceAllString(c, "")

	clearReg = regexp.MustCompile(`记住手机版网址：m.biqiuge.com`)
	c = clearReg.ReplaceAllString(c, "")
	return c
}
