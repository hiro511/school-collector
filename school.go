package main

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/shogo82148/go-mecab"
)

var (
	numWorkers              = 100
	wikiBase                *url.URL
	juniorHighSchoolListURL = "https://ja.wikipedia.org/wiki/%E6%97%A5%E6%9C%AC%E3%81%AE%E4%B8%AD%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7"
	elementarySchoolListURL = "https://ja.wikipedia.org/wiki/%E6%97%A5%E6%9C%AC%E3%81%AE%E5%B0%8F%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7"
	highSchoolURLs          = []string{
		"https://ja.wikipedia.org/wiki/%E5%8C%97%E6%B5%B7%E9%81%93%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 北海道
		"https://ja.wikipedia.org/wiki/%E9%9D%92%E6%A3%AE%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 青森
		"https://ja.wikipedia.org/wiki/%E5%B2%A9%E6%89%8B%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 秋田
		"https://ja.wikipedia.org/wiki/%E5%AE%AE%E5%9F%8E%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 宮城
		"https://ja.wikipedia.org/wiki/%E5%B1%B1%E5%BD%A2%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 山形
		"https://ja.wikipedia.org/wiki/%E7%A6%8F%E5%B3%B6%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 福島
		"https://ja.wikipedia.org/wiki/%E8%8C%A8%E5%9F%8E%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 茨木
		"https://ja.wikipedia.org/wiki/%E6%A0%83%E6%9C%A8%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 栃木
		"https://ja.wikipedia.org/wiki/%E7%BE%A4%E9%A6%AC%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 群馬
		"https://ja.wikipedia.org/wiki/%E5%9F%BC%E7%8E%89%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 埼玉
		"https://ja.wikipedia.org/wiki/%E5%8D%83%E8%91%89%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 千葉
		"https://ja.wikipedia.org/wiki/%E6%9D%B1%E4%BA%AC%E9%83%BD%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 東京
		"https://ja.wikipedia.org/wiki/%E7%A5%9E%E5%A5%88%E5%B7%9D%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7", // 神奈川
		"https://ja.wikipedia.org/wiki/%E6%96%B0%E6%BD%9F%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 新潟
		"https://ja.wikipedia.org/wiki/%E5%AF%8C%E5%B1%B1%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 富山
		"https://ja.wikipedia.org/wiki/%E7%9F%B3%E5%B7%9D%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 石川
		"https://ja.wikipedia.org/wiki/%E7%A6%8F%E4%BA%95%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 福井
		"https://ja.wikipedia.org/wiki/%E5%B1%B1%E6%A2%A8%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 山梨
		"https://ja.wikipedia.org/wiki/%E9%95%B7%E9%87%8E%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 長野
		"https://ja.wikipedia.org/wiki/%E5%B2%90%E9%98%9C%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 岐阜
		"https://ja.wikipedia.org/wiki/%E9%9D%99%E5%B2%A1%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 静岡
		"https://ja.wikipedia.org/wiki/%E6%84%9B%E7%9F%A5%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 愛知
		"https://ja.wikipedia.org/wiki/%E4%B8%89%E9%87%8D%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 三重
		"https://ja.wikipedia.org/wiki/%E6%BB%8B%E8%B3%80%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 滋賀
		"https://ja.wikipedia.org/wiki/%E4%BA%AC%E9%83%BD%E5%BA%9C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 京都
		"https://ja.wikipedia.org/wiki/%E5%A4%A7%E9%98%AA%E5%BA%9C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 大阪
		"https://ja.wikipedia.org/wiki/%E5%85%B5%E5%BA%AB%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 兵庫
		"https://ja.wikipedia.org/wiki/%E5%A5%88%E8%89%AF%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 奈良
		"https://ja.wikipedia.org/wiki/%E5%92%8C%E6%AD%8C%E5%B1%B1%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7", // 和歌山
		"https://ja.wikipedia.org/wiki/%E9%B3%A5%E5%8F%96%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 鳥取
		"https://ja.wikipedia.org/wiki/%E5%B3%B6%E6%A0%B9%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 島根
		"https://ja.wikipedia.org/wiki/%E5%B2%A1%E5%B1%B1%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 岡山
		"https://ja.wikipedia.org/wiki/%E5%BA%83%E5%B3%B6%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 広島
		"https://ja.wikipedia.org/wiki/%E5%B1%B1%E5%8F%A3%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 山口
		"https://ja.wikipedia.org/wiki/%E5%BE%B3%E5%B3%B6%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 徳島
		"https://ja.wikipedia.org/wiki/%E9%A6%99%E5%B7%9D%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 香川
		"https://ja.wikipedia.org/wiki/%E6%84%9B%E5%AA%9B%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 愛媛
		"https://ja.wikipedia.org/wiki/%E9%AB%98%E7%9F%A5%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 高知
		"https://ja.wikipedia.org/wiki/%E7%A6%8F%E5%B2%A1%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 福岡
		"https://ja.wikipedia.org/wiki/%E4%BD%90%E8%B3%80%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 佐賀
		"https://ja.wikipedia.org/wiki/%E9%95%B7%E5%B4%8E%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 長崎
		"https://ja.wikipedia.org/wiki/%E7%86%8A%E6%9C%AC%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 熊本
		"https://ja.wikipedia.org/wiki/%E5%A4%A7%E5%88%86%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 大分
		"https://ja.wikipedia.org/wiki/%E5%AE%AE%E5%B4%8E%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          // 宮崎
		"https://ja.wikipedia.org/wiki/%E9%B9%BF%E5%85%90%E5%B3%B6%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7", // 鹿児島
		"https://ja.wikipedia.org/wiki/%E6%B2%96%E7%B8%84%E7%9C%8C%E9%AB%98%E7%AD%89%E5%AD%A6%E6%A0%A1%E4%B8%80%E8%A6%A7",          //　沖縄
	}
	univURLs = []string{
		"https://ja.wikipedia.org/wiki/%E6%97%A5%E6%9C%AC%E3%81%AE%E5%A4%A7%E5%AD%A6%E4%B8%80%E8%A6%A7",
	}
)

type (
	schools struct {
		Suffix   string
		Results  []*result
		Last     string
		Finished bool
	}

	result struct {
		Name, Yomi, WikiURL string
	}
)

var (
	tasks = make(chan *result, 60000)
	wg    = &sync.WaitGroup{}
)

func (r result) Output() []string {
	return []string{r.Name, r.Yomi}
}

func init() {
	wiki, err := url.Parse("https://ja.wikipedia.org")
	if err != nil {
		panic(err)
	}
	wikiBase = wiki
}

func startWorker() {
	for i := 0; i < numWorkers; i++ {
		go do()
	}
}

func do() {

	for {
		select {
		case result, ok := <-tasks:
			if !ok {
				return
			}
			result.lookUpYomi()
			fmt.Printf("*")
			wg.Done()
		}
	}

}

func (sch *schools) Collect(URLs []string) {
	for i := range URLs {
		sch.collect(URLs[i])
	}
}

func (sch *schools) collect(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if sch.Finished {
			return
		}
		text := s.Text()
		if strings.HasSuffix(text, sch.Suffix) && !strings.ContainsRune(text, ' ') && text != sch.Suffix {
			r := &result{
				Name: text,
			}
			sch.Results = append(sch.Results, r)
			tasks <- r
			wg.Add(1)

			if r.Name == sch.Last {
				sch.Finished = true
			}
		}
	})
}

func resolveURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	if u.Path == "/w/index.php" {
		return ""
	}

	return wikiBase.ResolveReference(u).String()
}

var kanaConv = unicode.SpecialCase{
	unicode.CaseRange{
		0x30a1, // Lo: ァ
		0x30f3, // Hi: ン
		[unicode.MaxCase]rune{
			0,               // UpperCase では変換しない
			0x3041 - 0x30a1, // LowerCase でひらがなに変換
			0,               // TitleCase では変換しない
		},
	},
}

func (r *result) lookUpYomi() {
	tagger, err := mecab.New(map[string]string{"output-format-type": "yomi"})
	if err != nil {
		return
	}
	defer tagger.Destroy()

	yomi, err := tagger.Parse(r.Name)
	if err != nil {
		return
	}
	yomi = strings.TrimRight(yomi, "\n")

	r.Yomi = strings.ToLowerSpecial(kanaConv, yomi)
}

func (r *result) fetchYomi() {
	if r.WikiURL == "" {
		return
	}

	target := r.Name + "（"
	doc, err := goquery.NewDocument(r.WikiURL)
	if err != nil {
		return
	}

	text := doc.Text()
	start := strings.Index(text, target) + len(target)
	if start == -1 {
		return
	}
	length := strings.Index(text[start:], "）")
	if length == -1 || length > 200 {
		return
	}
	if strings.Contains(text[start:start+length], "、") {
		length = strings.Index(text[start:start+length], "、")
	}
	if strings.Contains(text[start:start+length], ", ") {
		length = strings.Index(text[start:start+length], ", ")
	}
	if strings.Contains(text[start:start+length], " ") {
		length = strings.Index(text[start:start+length], " ")
	}

	r.Yomi = text[start : start+length]
}

func schoolURLs(url, suffix string) []string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	results := make([]string, 0, 47)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.HasSuffix(text, suffix) {
			link, ok := s.Attr("href")
			if ok {
				results = append(results, resolveURL(link))
			}
		}
	})

	return results
}

func main() {
	startWorker()

	// highSchools := &schools{
	// 	Suffix: "高等学校",
	// }
	// highSchools.Collect(highSchoolURLs)
	//
	// juniorHighSchools := &schools{
	// 	Suffix: "中学校",
	// }
	// juniorHighSchoolURLs := schoolURLs(juniorHighSchoolListURL, "中学校一覧")
	// juniorHighSchools.Collect(juniorHighSchoolURLs)
	//
	// wg.Wait()
	// close(tasks)
	//
	// write(highSchools.Results, "highSchools.csv")
	// write(juniorHighSchools.Results, "juniorHighSchools.csv")

	// elementarySchools := &schools{
	// 	Suffix: "小学校",
	// }
	// elementarySchoolURLs := schoolURLs(elementarySchoolListURL, "小学校一覧")
	// elementarySchools.Collect(elementarySchoolURLs)
	// wg.Wait()
	// close(tasks)
	// write(elementarySchools.Results, "elementarySchools.csv")

	univ := &schools{
		Suffix: "大学",
		Last:   "和洋女子大学",
	}
	univ.Collect(univURLs)
	wg.Wait()
	close(tasks)
	write(univ.Results, "univ.csv")

}

func write(results []*result, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for i := range results {
		writer.Write(results[i].Output())
	}
	return nil
}
