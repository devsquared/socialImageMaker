package main

import (
	"bufio"
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	if err := cli(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func cli() error {
	fmt.Println("For best results, ensure that your images are 1200x628.")

	fmt.Println("Provide the image to set as the base. Put this in the 'backgroundImages' folder.")

	backgroundImageName, err := getInputText(); if err != nil {
		return err
	}

	fmt.Println("Provide the domain text located in lower left of image.")
	domainText, err := getInputText(); if err != nil {
		return err
	}

	fmt.Println("Provide a title for the main text in image.")
	title, err := getInputText(); if err != nil {
		return err
	}

	fmt.Println("Provide a path and name for your image. Can just provide a name with extension.")
	imageName, err := getInputText(); if err != nil {
		return err
	}

	if err := run(backgroundImageName, domainText, title, imageName); err != nil {
		return err
	}

	return nil
}

// create a basic CLI that will prompt for input to create a new image
// if we access the binary, run CLI. Otherwise, export the method for others to use
// here we will add parameters to customize the image
func run(bgImageName string, domainText string, titleText string, imageName string) error {
	dc := gg.NewContext(1200, 628)

	backgroundImage, err := gg.LoadImage(filepath.Join("./", "backgroundImages", bgImageName))
	if err != nil {
		return err
	}
	dc.DrawImage(backgroundImage, 0, 0)

	// NOTE!
	// add any transforms to the image here before the save

	// add a semi-transparent overlay here to add some contrast for text
	margin := 20.0
	x := margin
	y := margin
	w := float64(dc.Width()) - (2.0 * margin)
	h := float64(dc.Height()) - (2.0 * margin)
	dc.SetColor(color.RGBA{0, 0, 0, 204})
	dc.DrawRectangle(x, y, w, h)
	dc.Fill()

	// add text

	// first add to bottom righthand for domain name
	fontPath := filepath.Join("fonts", "FiraCode-Regular.ttf")
	if err := dc.LoadFontFace(fontPath, 80); err != nil {
		return err
	}
	dc.SetColor(color.White)
	s := domainText
	marginX := 50.0
	marginY := -10.0
	textWidth, textHeight := dc.MeasureString(s)
	x = float64(dc.Width()) - textWidth - marginX
	y = float64(dc.Height()) - textHeight - marginY
	dc.DrawString(s, x, y)

	// next add title
	title := titleText
	textShadowColor := color.Black
	textColor := color.White
	fontPath = filepath.Join("fonts", "FiraCode-Bold.ttf")
	if err := dc.LoadFontFace(fontPath, 90); err != nil {
		return err
	}
	textRightMargin := 60.0
	textTopMargin := 90.0
	x = textRightMargin
	y = textTopMargin
	maxWidth := float64(dc.Width()) - textRightMargin - textRightMargin
	dc.SetColor(textShadowColor)
	dc.DrawStringWrapped(title, x+1, y+1, 0, 0, maxWidth, 1.5, gg.AlignLeft)
	dc.SetColor(textColor)
	dc.DrawStringWrapped(title, x, y, 0, 0, maxWidth, 1.5, gg.AlignLeft)

	if err := dc.SavePNG(imageName); err != nil {
		return err
	}

	return nil
}

func getInputText() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n'); if err != nil {
		return "", err
	}

	text = strings.TrimSuffix(text, "\n")

	return text, nil
}

