package convertor

import (
	"fmt"
	"io"
	"os/exec"
)

type CropParameters struct {
	X      int
	Y      int
	Width  int
	Height int
}

type ResizeParameters struct {
	Width  int
	Height int
	Type   string
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

func (wp *WebP) SetCrop(crop CropParameters) *WebP {
	wp.parameters.crop = &crop
	return wp
}

func (wp *WebP) SetResize(resize ResizeParameters) *WebP {
	wp.parameters.resize = &resize

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
		args = append(args, "-crop", fmt.Sprintf("%d", wp.parameters.crop.X), fmt.Sprintf("%d", wp.parameters.crop.Y), fmt.Sprintf("%d", wp.parameters.crop.Width), fmt.Sprintf("%d", wp.parameters.crop.Height))
	}

	if wp.parameters.resize != nil {
		args = append(args, "-resize", fmt.Sprintf("%d", wp.parameters.resize.Width), fmt.Sprintf("%d", wp.parameters.resize.Height))
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
