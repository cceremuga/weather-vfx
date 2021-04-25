package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
	"github.com/icodealot/noaa"
	"golang.org/x/image/font"
)

type DailyForecast struct {
	Name     string
	HighTemp float64
	LowTemp  float64
}

func main() {
	forecast, err := sevenDayForecast("42.6526", "-73.7562")

	if err != nil {
		panic("Failed to retrieve forecast.")
	}

	if err = printForecastToImage(forecast); err != nil {
		panic("Failed to print forecast to file.")
	}
}

func sevenDayForecast(lat string, lng string) ([]DailyForecast, error) {
	forecast, err := noaa.Forecast(lat, lng)
	if err != nil {
		return nil, err
	}

	var sevenDayForecast = []DailyForecast{}
	var currentDay DailyForecast

	for i, period := range forecast.Periods {
		if i%2 == 0 {
			currentDay = DailyForecast{
				Name:     period.Name,
				HighTemp: period.Temperature,
			}
		} else {
			currentDay.LowTemp = period.Temperature
			sevenDayForecast = append(sevenDayForecast, currentDay)
		}
	}

	return sevenDayForecast, nil
}

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "luxisr.ttf", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 12, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", true, "white text on a black background")
)

func printForecastToImage(forecast []DailyForecast) error {
	for _, day := range forecast {
		fmt.Printf("%-20s %.0f / %.0f\n", day.Name, day.HighTemp, day.LowTemp)
	}

	flag.Parse()

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		return err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	// Initialize the context.
	fg, bg := image.Black, image.White
	ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff}
	if *wonb {
		fg, bg = image.White, image.Black
		ruler = color.RGBA{0x22, 0x22, 0x22, 0xff}
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the guidelines.
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	text := []string{"this", "is", "a", "test"}

	// Draw the text.
	pt := freetype.Pt(10, 10+int(c.PointToFixed(*size)>>6))
	for _, s := range text {
		_, err = c.DrawString(s, pt)
		if err != nil {
			return err
		}
		pt.Y += c.PointToFixed(*size * *spacing)
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create("out/out.png")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")

	return nil
}
