package content

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
)

type cont struct {
	Title   string
	Content string
	Url     string
}

// GetLastTitleList 获取最近更新列表
func GetLastTitleList(url string, num int) []cont {
	//url := "https://m.biqiuge.com/"
	var ret []cont
	doc, _ := get(url)
	doc.Find(".books .listpage select option").Each(func(i int, s *goquery.Selection) {
		values, _ := s.Attr("value")
		c := cont{Url: values}
		ret = append(ret, c)
	})

	//取最后一页上的所有章节名
	lastC := ret[len(ret)-1]
	titleListDoc, _ := get(url + lastC.Url)
	titleLists := getTitleList(titleListDoc)

	//取最后一页的内容太少，再取前一页，防止一次更新太多时，会缺少
	if len(titleLists) < num {
		secondC := ret[len(ret)-2]
		c, _ := get(url + secondC.Url)
		titleLists = append(getTitleList(c))
	}

	fmt.Printf("%+v\n", titleLists)
	return ret
}

func GetCont(u string) string {
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
			fmt.Print("%v\n", err)
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
		return nil, err
	}
	defer res.Body.Close()

	utfBody, err := iconv.NewReader(res.Body, "gbk", "utf-8")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func getTitleList(titleListDoc *goquery.Document) []cont {
	var titleLists []cont
	titleListDoc.Find(".books .book_last").Each(func(i int, s *goquery.Selection) {
		if i == 1 {
			s.Find("dl dd a").Each(func(ii int, ss *goquery.Selection) {
				urlStr, eUrlStr := ss.Attr("href")
				if eUrlStr {
					t := cont{Url: urlStr, Title: ss.Text()}
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
