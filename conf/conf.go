package conf

import (
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"strings"
	"time"
)

var Genre string
var Rating float64
var UserAgent string
var Cookies string

func Init() {
	config.WithOptions(config.ParseEnv)

	// add Decoder and Encoder
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("setting.yml")
	if err != nil {
		panic(err)
	}

	genre := strings.TrimSpace(config.String("genre"))
	rating := config.Float("rating")
	userAgent := strings.TrimSpace(config.String("userAgent"))
	cookies := strings.TrimSpace(config.String("cookies"))

	if genre == "" {
		panic("Genre in setting.yml can not be empty!!!")
	}

	if userAgent == "" {
		panic("UserAgent in setting.yml can not be empty!!!")
	}

	if cookies == "" {
		panic("Cookies in setting.yml can not be empty!!!")
	}

	fmt.Printf("Setting load succeed, Start in 5s:\n - Genre: %s \n - Rating: %f\n\n", genre, rating)

	Genre = genre
	Rating = rating
	UserAgent = userAgent
	Cookies = cookies

	time.Sleep(5 * time.Second)
}
