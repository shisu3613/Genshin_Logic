package game

import (
	"fmt"
	"regexp"
	"time"
)

//单例模式

var manageBanWord *ManageBanWord

type ManageBanWord struct {
	BanWordBase  []string //配置生成
	BanWordExtra []string //更新
}

func GetManageBanWord() *ManageBanWord {
	if manageBanWord == nil {
		manageBanWord = new(ManageBanWord)
		manageBanWord.BanWordBase = []string{"外挂", "外挂工具"}
		manageBanWord.BanWordExtra = []string{"原神", "外挂工具"}
	}
	return manageBanWord
}

func (mbw *ManageBanWord) IsBanWord(txt string) bool {
	for _, v := range mbw.BanWordBase {
		match, _ := regexp.MatchString(v, txt)
		if match {
			fmt.Println(match, v)
			return match
		}
	}

	for _, v := range mbw.BanWordExtra {
		match, _ := regexp.MatchString(v, txt)
		if match {
			fmt.Println(match, v)
			return match
		}
	}
	return false
}

func (mbw *ManageBanWord) Run() {
	tickerIn := time.NewTicker(time.Second * 1)
	//tickerOut:= time.NewTicker(time.Second * 5)
	for {
		select {
		case <-tickerIn.C:
			if time.Now().Unix()%10 == 0 {
				fmt.Println("更新词库")
			} else {
				//fmt.Println("待机")
			}
		}
	}
}
