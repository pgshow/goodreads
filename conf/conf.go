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

	if genre == "" {
		panic("Genre in setting.yml can not be empty!!!")
	}

	fmt.Printf("Setting load succeed, Start in 5s:\n - Genre: %s \n - Rating: %f\n\n", genre, rating)

	Genre = genre
	Rating = rating

	time.Sleep(5 * time.Second)
}
