package saveExcel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"goodreads/conf"
	"goodreads/model"
)

func Save(books []model.Book)  {
	f := excelize.NewFile()
	// 设置表格头
	f.SetCellValue("Sheet1", "A1", "Title")
	f.SetCellValue("Sheet1", "B1", "Author")
	f.SetCellValue("Sheet1", "C1", "Published_Year")
	f.SetCellValue("Sheet1", "D1", "Genre")
	f.SetCellValue("Sheet1", "E1", "Average_Ratings")
	f.SetCellValue("Sheet1", "F1", "Number_Ratings")
	f.SetCellValue("Sheet1", "G1", "Image_URL")

	// 循环写入数据
	line := 1
	for _, b := range books {
		line++
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", line), b.Title)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", line), b.Author)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", line), b.Published_Year)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", line), b.Genre)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", line), b.Average_Ratings)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", line), b.Number_Ratings)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", line), b.Image_URL)
	}

	// 根据指定路径保存文件
	if err := f.SaveAs(fmt.Sprintf("./%s.xlsx", conf.Genre)); err != nil {
		fmt.Println(err)
	}
}