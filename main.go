package main

import (
	"fiction/db"
	"fmt"
)

func main() {
	// content.GetLastTitleList(5)
	db.SetTable("fuckt")
	//i, err := db.SaveCont("fuck wangt", "fuck fuck fuck")
	//s, err := db.GetContByMd5("32bf0e6fcff51e53bd74e70ba1d622b2")
	s, err := db.GetLasterFive()
	fmt.Println(s)
	fmt.Println(err)
}
