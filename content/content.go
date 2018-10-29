package content

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
)

type cont struct {
	Title   string
	Content string
	Url     string
}

// GetLastTitleList 获取最近更新列表
func GetLastTitleList(num int) []cont {
	url := "https://m.biqiuge.com/"
	var ret []cont
	doc := get(url + "book_4772/")
	doc.Find(".books .listpage select option").Each(func(i int, s *goquery.Selection) {
		values, _ := s.Attr("value")
		c := cont{Url: values}
		ret = append(ret, c)
	})

	//取最后一页上的所有章节名
	lastC := ret[len(ret)-1]
	titleListDoc := get(url + lastC.Url)
	titleLists := getTitleList(titleListDoc)

	//取最后一页的内容太少，再取前一页，防止一次更新太多时，会缺少
	if len(titleLists) < num {
		secondC := ret[len(ret)-2]
		titleLists = append(getTitleList(get(url + secondC.Url)))
	}

	//fmt.Printf("%+v\n", titleLists)
	return ret
}

func get(url string) *goquery.Document {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	req.Header.Set("Referer", url)
	req.Header.Set("User-Agent", " Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1")
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s\n", "Error Network")
	}
	defer res.Body.Close()

	utfBody, err := iconv.NewReader(res.Body, "gbk", "utf-8")
	if err != nil {
		fmt.Printf("%s\n", "Error Encoding")
		os.Exit(-1)
	}
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		fmt.Printf("%s", "Error decode body")
		os.Exit(-1)
	}
	return doc
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
