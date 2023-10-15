package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	var lDiameter float64
	var rDiameter float64
	var dateStr string

	flag.Float64Var(&lDiameter, "left-diam", 17.0, "Circle Diameter for Left Hand to Draw")
	flag.Float64Var(&rDiameter, "right-diam", 17.0, "Circle Diameter for Right Hand to Draw")
	flag.StringVar(&dateStr, "date", "today", "Custom date in YYYY-MM-DD format")

	flag.Parse()

	date := determineDate(dateStr)

	err := generatePDF(lDiameter, rDiameter, date)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(`Generated PDF with
  date                       %s
  left hand circle diameter  %.1f mm
  right hand circle diameter %.1f mm`,
		date,
		lDiameter,
		rDiameter)
}

func determineDate(dateStr string) string {
	if dateStr != "today" {
		return dateStr
	}
	return time.Now().Format("2006-01-02")
}

func generatePDF(lDiameter float64, rDiameter float64, date string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)

	for i := 1; i <= 4; i++ {
		var diameter float64
		if i%2 == 0 {
			diameter = rDiameter
		} else {
			diameter = lDiameter
		}
		pdf.AddPage()
		drawHeader(pdf, date, i, diameter)
		drawGrid(pdf, diameter)
	}

	return pdf.OutputFileAndClose("output.pdf")
}

func drawHeader(pdf *gofpdf.Fpdf, date string, pageNumber int, diameter float64) {
	var patch string
	var pen string
	switch pageNumber {
	case 1:
		patch = "Left"
		pen = "Left"
	case 2:
		patch = "Left"
		pen = "Right"
	case 3:
		patch = "Right"
		pen = "Left"
	case 4:
		patch = "Right"
		pen = "Right"
	default:
		return
	}
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, fmt.Sprintf("%s Eye Patched, Pen in %s Hand | %s | Diameter: %.1fmm", patch, pen, date, diameter))
	pdf.Ln(20)
}

func drawGrid(pdf *gofpdf.Fpdf, diameter float64) {
	margin := 10.0
	radius := diameter / 2.0
	for x := margin + radius; x <= 210-margin-radius; x += diameter {
		for y := 30.0 + radius; y <= 297.0-margin-radius; y += diameter {
			pdf.Circle(x, y, radius, "D")
		}
	}
}
