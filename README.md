# Go MNIST

## Installing
```bash 
go get github.com/sno6/mnist
```

## Getting started

```go
// Read images & labels from local ubyte files.
r, err := mnist.New(&mnist.Files{
    TrainingImagesLoc: imagesLoc,
    TrainingLabelsLoc: labelsLoc,
})

...

// Use the reader to access files & labels.
img := r.TrainingImages.GetImage(1)
label := r.TrainingLabels.GetLabel(1)

...
```

## Example
Working code example: `/example`
