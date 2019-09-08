package mnist

import (
	"encoding/binary"
	"image"
	"image/color"
	"io"
	"os"

	"github.com/pkg/errors"
)

const (
	ImageMagic int32 = 2051
	LabelMagic int32 = 2049
)

var errMagicMismatch = errors.New("mnist: magic mismatch")

type Files struct {
	TrainingImagesLoc string
	TrainingLabelsLoc string
	TestingImagesLoc  string
	TestingLabelsLoc  string
}

type Reader struct {
	TrainingImages, TestingImages *ImageSet
	TrainingLabels, TestingLabels *LabelSet
}

type ImageSet struct {
	Count      int
	Rows, Cols int
	Images     []uint8
}

type LabelSet struct {
	Count  int
	Labels []uint8
}

type binLabelCfg struct {
	Count int32
}

type binImageCfg struct {
	Count      int32
	Rows, Cols int32
}

func New(locs *Files) (*Reader, error) {
	r := &Reader{}
	err := r.process(locs)
	return r, err
}

func (r *Reader) process(locs *Files) error {
	var err error
	if locs.TrainingImagesLoc != "" {
		r.TrainingImages, err = r.readImageSet(locs.TrainingImagesLoc)
		if err != nil {
			return err
		}
	}
	if locs.TrainingLabelsLoc != "" {
		r.TrainingLabels, err = r.readLabelSet(locs.TrainingLabelsLoc)
		if err != nil {
			return err
		}
	}
	if locs.TestingImagesLoc != "" {
		r.TestingImages, err = r.readImageSet(locs.TestingImagesLoc)
		if err != nil {
			return err
		}
	}
	if locs.TestingLabelsLoc != "" {
		r.TestingLabels, err = r.readLabelSet(locs.TestingLabelsLoc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reader) readLabelSet(loc string) (*LabelSet, error) {
	fr, err := newReaderFromFile(loc, LabelMagic)
	if err != nil {
		return nil, errors.Wrap(err, "readLabelSet: error reading label file")
	}

	var bc binLabelCfg
	if err := binary.Read(fr, binary.BigEndian, &bc); err != nil {
		return nil, errors.Wrap(err, "readLabelset: error reading label config")
	}

	d := make([]uint8, int(bc.Count))
	if err := binary.Read(fr, binary.BigEndian, d); err != nil {
		return nil, errors.Wrap(err, "readLabelset: error reading label data")
	}

	return &LabelSet{
		Count:  int(bc.Count),
		Labels: d,
	}, nil
}

func (r *Reader) readImageSet(loc string) (*ImageSet, error) {
	fr, err := newReaderFromFile(loc, ImageMagic)
	if err != nil {
		return nil, errors.Wrap(err, "readImageSet: error reading image file")
	}

	var bc binImageCfg
	if err := binary.Read(fr, binary.BigEndian, &bc); err != nil {
		return nil, errors.Wrap(err, "readImageSet: error reading image config")
	}

	d := make([]uint8, int(bc.Count*bc.Rows*bc.Cols))
	if err := binary.Read(fr, binary.BigEndian, d); err != nil {
		return nil, errors.Wrap(err, "readImageSet: error reading image data")
	}

	return &ImageSet{
		Count:  int(bc.Count),
		Rows:   int(bc.Rows),
		Cols:   int(bc.Cols),
		Images: d,
	}, nil
}

func (l *LabelSet) GetLabel(offset int) uint8 {
	if offset > l.Count {
		panic("GetLabel: offset value greater than count")
	}
	return l.Labels[offset]
}

func (img *ImageSet) GetImage(offset int) [][]uint8 {
	if offset > img.Count {
		panic("GetImage: offset value greater than count")
	}

	m := make([][]uint8, img.Rows)
	for i := range m {
		m[i] = make([]uint8, img.Cols)
	}

	r, c := len(m), len(m[0])
	imgStart := offset * r * c
	for y := 0; y < r; y++ {
		for x := 0; x < c; x++ {
			valOffset := (y * r) + x
			m[x][y] = img.Images[imgStart+valOffset]
		}
	}
	return m
}

func (img *ImageSet) GetFlatImage(offset int) []uint8 {
	if offset > img.Count {
		panic("GetFlatImage: offset value greater than count")
	}

	sz := img.Rows * img.Cols
	return img.Images[offset*sz : (offset*sz)+sz]

}

func ToGrayScale(m [][]uint8) image.Image {
	if len(m) == 0 || m == nil {
		return nil
	}

	r, c := len(m), len(m[0])
	img := image.NewGray(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{r, c}})
	for y := 0; y < r; y++ {
		for x := 0; x < c; x++ {
			img.Set(x, y, color.Gray{Y: m[x][y]})
		}
	}
	return img
}

func newReaderFromFile(loc string, magic int32) (io.Reader, error) {
	r, err := os.Open(loc)
	if err != nil {
		return nil, errors.Wrap(err, "newReaderFromFile: error opening file")
	}
	if err := checkMagic(r, magic); err != nil {
		return nil, err
	}
	return r, nil
}

func checkMagic(r io.Reader, truth int32) error {
	var magic int32
	if err := binary.Read(r, binary.BigEndian, &magic); err != nil {
		return err
	}
	if magic != truth {
		return errMagicMismatch
	}
	return nil
}
