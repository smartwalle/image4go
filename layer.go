package image4go

import (
	"bufio"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"reflect"
)

// LayerAlignment 水平对齐
type LayerAlignment int

// LayerVerticalAlignment 垂直对齐
type LayerVerticalAlignment int

const (
	LayerAlignmentDefault LayerAlignment = iota
	LayerAlignmentLeft
	LayerAlignmentCenter
	LayerAlignmentRight
)

const (
	LayerVerticalAlignmentDefault LayerVerticalAlignment = iota
	LayerVerticalAlignmentTop
	LayerVerticalAlignmentMiddle
	LayerVerticalAlignmentBottom
)

type Layer interface {
	Render() image.Image

	Rect() image.Rectangle

	SetAlignment(alignment LayerAlignment)

	Alignment() LayerAlignment

	SetVerticalAlignment(alignment LayerVerticalAlignment)

	VerticalAlignment() LayerVerticalAlignment
}

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) Point {
	return Point{X: x, Y: y}
}

type Size struct {
	Width  int
	Height int
}

func NewSize(width, height int) Size {
	return Size{Width: width, Height: height}
}

func isNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

func calcRect(pRect, sRect image.Rectangle, alignment LayerAlignment, verticalAlignment LayerVerticalAlignment) (rect image.Rectangle) {
	var pWidth = pRect.Max.X - pRect.Min.X
	var sWidth = sRect.Max.X - sRect.Min.X

	switch alignment {
	case LayerAlignmentDefault:
		rect.Min.X = sRect.Min.X
		rect.Max.X = sRect.Max.X
	case LayerAlignmentLeft:
		rect.Min.X = 0
		rect.Max.X = sWidth
	case LayerAlignmentCenter:
		var w = pWidth - sWidth
		rect.Min.X = w / 2
		rect.Max.X = rect.Min.X + sWidth
	case LayerAlignmentRight:
		rect.Min.X = pWidth - sWidth
		rect.Max.X = rect.Min.X + sWidth
	default:
		rect.Min.X = sRect.Min.X
		rect.Max.X = sRect.Max.X
	}

	var pHeight = pRect.Max.Y - pRect.Min.Y
	var sHeight = sRect.Max.Y - sRect.Min.Y

	switch verticalAlignment {
	case LayerVerticalAlignmentDefault:
		rect.Min.Y = sRect.Min.Y
		rect.Max.Y = sRect.Max.Y
	case LayerVerticalAlignmentTop:
		rect.Min.Y = 0
		rect.Max.Y = sHeight
	case LayerVerticalAlignmentMiddle:
		var h = pHeight - sHeight
		rect.Min.Y = h / 2
		rect.Max.Y = rect.Min.Y + sHeight
	case LayerVerticalAlignmentBottom:
		rect.Min.Y = pHeight - sHeight
		rect.Max.Y = rect.Min.Y + sHeight
	default:
		rect.Min.Y = sRect.Min.Y
		rect.Max.Y = sRect.Max.Y
	}
	return rect
}

func WriteToPNG(l Layer, file string) (err error) {
	nFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer nFile.Close()

	b := bufio.NewWriter(nFile)

	if err = png.Encode(nFile, l.Render()); err != nil {
		return err
	}
	return b.Flush()
}

func WriteToJPEG(l Layer, file string, quality int) (err error) {
	nFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer nFile.Close()

	b := bufio.NewWriter(nFile)

	if err = jpeg.Encode(nFile, l.Render(), &jpeg.Options{Quality: quality}); err != nil {
		return err
	}
	return b.Flush()
}
