package save

import (
	"encoding/csv"
	"fmt"
	"goodreads/conf"
	"goodreads/model"
	"io"
	"log"
	"os"
	"reflect"
)

func Save(books []model.Book) {
	createNew(fmt.Sprintf("./%s.csv", conf.Genre))

	for _, book := range books{
		add(&book, fmt.Sprintf("./%s.csv", conf.Genre))
	}
}

func createNew(path string) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("can not create file: ", err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF")

	writer := csv.NewWriter(f)
	defer writer.Flush()

	//将爬取信息写入csv文件
	writer.Write([]string{
		"Title",
		"Author",
		"Published_Year",
		"Genre",
		"Average_Ratings",
		"Number_Ratings",
		"Image_URL",
	})
}

func add(item *model.Book, path string) {
	if !reflect.ValueOf(item).IsValid() {
		return
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("can not open file %s, err is %+v", path, err)
	}
	defer f.Close()
	f.Seek(0, io.SeekEnd)

	w := csv.NewWriter(f)
	//设置属性
	w.UseCRLF = true
	row := struct2strings(item)
	err = w.Write(row)
	if err != nil {
		log.Fatalf("can't write to %s, err is %+v", path, err)
	}
	//这里必须刷新，才能将数据写入文件。
	w.Flush()
}


// struct 转 []string
func struct2strings(item *model.Book) (result []string) {
	v := reflect.ValueOf(*item)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		result = append(result, f.String())
	}
	return
}