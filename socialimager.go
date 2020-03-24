package main

import (
	"bufio"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image/color"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func main() {

	if err := cli(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func cli() error {
	// setup a ctrl-c interrupt to cancel and exit the program
	cancelled := setupCloseHandler()

	// TODO: Let's add the ability to have a "default" flag
	// This will pull from some config.yaml that will then fill in the blanks except for title

	fmt.Println("Welcome to the Social Image Maker!")
	fmt.Println("------------------------------------")
	fmt.Println("For best results, ensure that your images are 1200x628.")
	fmt.Println("Provide the image to set as the base. Put this in the 'backgroundImages' folder.")

	backgroundImageName, err := getInputText()
	if err != nil {
		return err
	}

	fmt.Println("Enter a size for the image.")
	fmt.Println("First, width:")
	imageWidthStr, err := getInputText()
	if err != nil {
		return err
	}
	fmt.Println("Then, height:")
	imageHeightStr, err := getInputText()
	if err != nil {
		return err
	}

	fmt.Println("Provide the domain text located in lower left of image.")
	domainText, err := getInputText()
	if err != nil {
		return err
	}

	fmt.Println("What font size should I make this domain text?")
	domainFontSizeStr, err := getInputText()
	if err != nil {
		return err
	}

	fmt.Println("Provide a title for the main text in image.")
	title, err := getInputText()
	if err != nil {
		return err
	}

	fmt.Println("What font size should I make this domain text?")
	titleFontSizeStr, err := getInputText()
	if err != nil {
		return err
	}

	fmt.Println("Provide a path and name for your image. Can just provide a name with extension.")
	imageName, err := getInputText()
	if err != nil {
		return err
	}

	if !cancelled {
		imageWidth, err := strconv.Atoi(imageWidthStr)
		if err != nil {
			return err
		}

		imageHeight, err := strconv.Atoi(imageHeightStr)
		if err != nil {
			return err
		}

		domainFontSizeInt, err := strconv.Atoi(domainFontSizeStr)
		if err != nil {
			return err
		}
		domainFontSize := float64(domainFontSizeInt)

		titleFontSizeInt, err := strconv.Atoi(titleFontSizeStr)
		if err != nil {
			return err
		}
		titleFontSize := float64(titleFontSizeInt)

		if err := run(backgroundImageName, domainText, title, imageName, imageWidth, imageHeight,
			domainFontSize, titleFontSize); err != nil {
			return err
		}

		fmt.Println("All done! Check for your new image at the path you provided!")
	}

	return nil
}

// create a basic CLI that will prompt for input to create a new image
// if we access the binary, run CLI. Otherwise, export the method for others to use
// here we will add parameters to customize the image

// the original measurements were 1200x628 with domain font size of 80 and title font size of 90
func run(bgImageName string, domainText string, titleText string, imageName string, imageWidth int, imageHeight int,
	domainFontSize float64, titleFontSize float64) error {
	dc := gg.NewContext(imageWidth, imageHeight)

	backgroundImage, err := gg.LoadImage(filepath.Join("./", "backgroundImages", bgImageName))
	if err != nil {
		return err
	}
	resizedImage := imaging.Resize(backgroundImage, imageWidth, imageHeight, imaging.Lanczos)

	dc.DrawImage(resizedImage, 0, 0)

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
	if err := dc.LoadFontFace(fontPath, domainFontSize); err != nil {
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
	if err := dc.LoadFontFace(fontPath, titleFontSize); err != nil {
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
	// prompt for input
	fmt.Print(">> ")

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	text = strings.TrimSuffix(text, "\n")

	return text, nil
}

// NOTE: This close handler does not seem to work in situations where a IDE handles the running of the app
// It, however, works as intended in a terminal.
func setupCloseHandler() bool {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() bool {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
		return true
	}()

	//base case
	return false
}
