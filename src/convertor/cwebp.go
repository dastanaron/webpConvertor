package convertor

import (
	"fmt"
	"io"
	"os/exec"
)

type CropParameters struct {
	x      int
	y      int
	width  int
	heigth int
}

type ResizeParameters struct {
	width  int
	height int
}

type ConvertParams struct {
	crop    *CropParameters
	resize  *ResizeParameters
	quality int
}

type WebP struct {
	input      io.Reader
	output     io.Writer
	binPath    string
	parameters *ConvertParams
}

func NewCWebP() *WebP {
	webp := &WebP{
		parameters: &ConvertParams{
			quality: 100,
		},
	}

	return webp
}

func (wp *WebP) Input(input io.Reader) *WebP {
	wp.input = input
	return wp
}

func (wp *WebP) Output(output io.Writer) *WebP {
	wp.output = output
	return wp
}

func (wp *WebP) SetBinPath(path string) *WebP {
	wp.binPath = path
	return wp
}

func (wp *WebP) SetCrop(x, y, width, heigth int) *WebP {
	wp.parameters.crop = &CropParameters{
		x:      x,
		y:      y,
		width:  width,
		heigth: heigth,
	}

	return wp
}

func (wp *WebP) SetResize(width, heigth int) *WebP {
	wp.parameters.resize = &ResizeParameters{
		width:  width,
		height: heigth,
	}

	return wp
}

func (wp *WebP) SetQuality(quality int) *WebP {
	wp.parameters.quality = quality
	return wp
}

func (wp *WebP) Run() error {

	var args []string

	command := fmt.Sprintf("%s/cwebp", wp.binPath)

	if wp.parameters.crop != nil {
		args = append(args, "-crop", fmt.Sprintf("%d", wp.parameters.crop.x), fmt.Sprintf("%d", wp.parameters.crop.y), fmt.Sprintf("%d", wp.parameters.crop.width), fmt.Sprintf("%d", wp.parameters.crop.heigth))
	}

	if wp.parameters.resize != nil {
		args = append(args, "-resize", fmt.Sprintf("%d", wp.parameters.resize.width), fmt.Sprintf("%d", wp.parameters.resize.height))
	}

	args = append(args, "-q", fmt.Sprintf("%d", wp.parameters.quality))

	args = append(args, "-o", "-")

	args = append(args, "--", "-")

	cmd := exec.Command(command, args...)

	cmd.Stdin = wp.input
	cmd.Stdout = wp.output

	err := cmd.Start()

	if err != nil {
		return err
	}

	err = cmd.Wait()

	return err
}
