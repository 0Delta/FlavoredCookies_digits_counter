package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Entry struct {
	name  string
	value int
}
type List []Entry

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	if l[i].value == l[j].value {
		return (l[i].name < l[j].name)
	} else {
		return (l[i].value < l[j].value)
	}
}

func main() {
	get("https://w.atwiki.jp/cookieclickerjpn/pages/46.html", "./temp.html")
	ret := screpUpList("./temp.html")
	sort.Sort(ret)
	for _, v := range ret {
		fmt.Println(v.value, v.name)
	}
}

func get(url string, fname string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Print("url scarapping failed")
	}
	res, err := doc.Find("body").Html()
	if err != nil {
		fmt.Print("dom get failed")
	}
	ioutil.WriteFile(fname, []byte(res), os.ModePerm)
}

func screpUpList(fname string) List {
	fileInfos, _ := ioutil.ReadFile(fname)
	stringReader := strings.NewReader(string(fileInfos))
	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		fmt.Print("url scarapping failed")
	}
	ret := List{}

	fmt.Println("test2")
	doc.Find("#wikibody > table > tbody > tr").Each(func(_ int, s *goquery.Selection) {
		str := s.Find("td:nth-child(1)").Text()
		str2, _ := s.Find("td:nth-child(4) > span").Attr("title")
		// str, _ := s.Attr("title")

		str2 = strings.Replace(str2, ",", "", -1)
		str2i := utf8.RuneCountInString(str2)
		ret = append(ret, Entry{str, str2i})
	})
	return ret
}

func getCharaPage(url string, subdir string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Print("url scarapping failed")
	}
	res := doc.Find("h2").Text()
	str, _, err := transform.String(japanese.EUCJP.NewDecoder(), res)
	if err != nil {
		fmt.Print("dom get failed")
	}
	str = strings.TrimSpace(str)
	strs := strings.Split(str, "\n")
	fmt.Println(strs[0])

	res = doc.Find("#page-body-inner > div.user-area").Text()
	str, _, err = transform.String(japanese.EUCJP.NewDecoder(), res)
	if err != nil {
		fmt.Print("dom get failed")
	}
	//	fmt.Println(str)
	fname := "./pages/" + subdir + "/" + strs[0]
	ioutil.WriteFile(fname, []byte(str), os.ModePerm)
}
