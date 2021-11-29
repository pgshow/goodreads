package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"goodreads/conf"
	"goodreads/model"
	"goodreads/saveExcel"
	"goodreads/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	conf.Init()

	var bookList []model.Book
	var i = 0
	for true {
		time.Sleep(20 * time.Second)
		i += 1
		url := fmt.Sprintf("https://www.goodreads.com/shelf/show/%s?page=%d", conf.Genre, i)
		println("scrape: " + url)

		body, _ := util.Fetch(url)

		if strings.Contains(body, "Back to the Goodreads homepage") {
			println("The page is over, stop scraping")
			break
		}

		if strings.Contains(body, "Showing 0-0 of 0") {
			println("Nothing on this page")
			break
		}

		doc, err := htmlquery.Parse(strings.NewReader(body))
		if err != nil {
			return
		}

		bookItems := htmlquery.Find(doc, "//div[@class='leftContainer']/div[@class='elementList']")

		if bookItems == nil {
			return
		}

		k := 0
		for _, book := range bookItems {
			k += 1

			var bookOjb model.Book

			// 标题
			titleTmp := htmlquery.FindOne(book, ".//a[@class='bookTitle']")
			if titleTmp == nil {
				println("No title for this item, Pass it")
				continue
			}

			bookOjb.Title = htmlquery.InnerText(titleTmp)

			// 图片
			imgTmp := htmlquery.FindOne(book, ".//img")
			if imgTmp != nil {
				thumbnail := htmlquery.SelectAttr(imgTmp, "src")                         // 略缩图
				bookOjb.Image_URL = strings.Replace(thumbnail, "._SX50_SY75_", "", 1)    // 原图
				bookOjb.Image_URL = strings.Replace(bookOjb.Image_URL, "._SY75_", "", 1) // 原图
				bookOjb.Image_URL = strings.Replace(bookOjb.Image_URL, "._SX50_", "", 1) // 原图
			}

			// 作者
			authorTmp := htmlquery.FindOne(book, ".//span[@itemprop='name']")
			if authorTmp != nil {
				bookOjb.Author = htmlquery.InnerText(authorTmp)
			}

			// 类别
			genreTmp := htmlquery.FindOne(book, ".//em")
			if genreTmp != nil {
				bookOjb.Genre = htmlquery.InnerText(genreTmp)
			}

			//if bookOjb.Title == "Autobiography of a Yogi (Paperback)" {
			//	println(111)
			//}

			// 发行年份, 平均评分, 评分人数
			otherTmp := htmlquery.FindOne(book, ".//span[@class='greyText smallText']")
			if otherTmp != nil {
				textTmp := htmlquery.InnerText(otherTmp)
				textTmp = strings.Replace(textTmp, "\n", "", -1)

				if textTmp != "" {
					regexpRule := regexp.MustCompile(`avg rating ([\d\.]+)[ —]+([\d,]+) ratings[ —]+published (\d{0,4})`)
					matchTmp := regexpRule.FindStringSubmatch(textTmp)

					if tmp, err := strconv.ParseFloat(matchTmp[1], 64); err == nil {
						bookOjb.Average_Ratings = tmp
					}

					if tmp, err := strconv.Atoi(strings.ReplaceAll(matchTmp[2], ",", "")); err == nil {
						bookOjb.Number_Ratings = tmp
					}

					bookOjb.Published_Year = matchTmp[3]
				}
			}

			// 抛弃低评分的书籍
			if err != nil || bookOjb.Average_Ratings <= conf.Rating {
				fmt.Println(k, bookOjb.Title, "- rating is below", conf.Rating)
				continue
			}

			fmt.Println(k, bookOjb.Title, "-", bookOjb.Author)

			bookList = append(bookList, bookOjb)
		}
	}

	saveExcel.Save(bookList)
}
