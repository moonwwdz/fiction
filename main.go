package main

import (
	"log"
	"sync"
	"time"

	"github.com/moonwwdz/fiction/config"
	"github.com/moonwwdz/fiction/content"
	"github.com/moonwwdz/fiction/db"
	"github.com/moonwwdz/fiction/util"
)

var wg sync.WaitGroup

func main() {

	confData, err := config.GetConf()
	if err != nil {
		log.Println("读取配置信息失败，需要重新生成")
	}

	log.Println(content.GetCont(""))

	if len(confData.Conf) > 0 {
		// wg.Add(len(confData.Conf))
		for _, sConf := range confData.Conf {
			// go deal(sConf)
			db.SetTable(sConf.TableName)
			time.Sleep(5 * time.Second)
			titleData := content.GetLastTitleList(sConf.Url, 10)
			for _, sTitle := range titleData {
				log.Println(sTitle.Title)
				titleMd5 := util.GetMD5Hash(sTitle.Title)
				dbCont, err := db.GetContByMd5(titleMd5)
				if err != nil || dbCont == "" {
					sCont := content.GetCont(sTitle.Url)
					_, err := db.SaveCont(sTitle.Title, sCont)
					if err == nil {
						log.Println("Save new caption of :" + sTitle.Title)
					}
				}
			}
		}
		// wg.Wait()
	}
}

func deal(fConf config.Fconf) {
	defer wg.Done()
	titleCont := content.GetLastTitleList(fConf.Url, 10)
	log.Printf("%v", titleCont)
}
