package assets

import (
	"embed"
	"io"

	_ "image/png"

	resource "github.com/quasilyte/ebitengine-resource"
)

const (
	ImgNone resource.ImageID = iota
	ImgGround2D
	ImgGroundIso
	ImgBlueCube2D
	ImgBlueCubeIso
	ImgBlueHalfCube2D
	ImgBlueHalfCubeIso
	ImgRedCube2D
	ImgRedCubeIso
	ImgRedHalfCube2D
	ImgRedHalfCubeIso
	ImgYellowCube2D
	ImgYellowCubeIso
	ImgYellowHalfCube2D
	ImgYellowHalfCubeIso
)

func RegisterImageResources(loader *resource.Loader) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImgGround2D:          {Path: "ground-2d.png"},
		ImgGroundIso:         {Path: "ground-iso.png"},
		ImgBlueCube2D:        {Path: "blue-2d-cube.png"},
		ImgBlueCubeIso:       {Path: "blue-iso-cube.png"},
		ImgBlueHalfCube2D:    {Path: "blue-2d-half-cube.png"},
		ImgBlueHalfCubeIso:   {Path: "blue-iso-half-cube.png"},
		ImgRedCube2D:         {Path: "red-2d-cube.png"},
		ImgRedCubeIso:        {Path: "red-iso-cube.png"},
		ImgRedHalfCube2D:     {Path: "red-2d-half-cube.png"},
		ImgRedHalfCubeIso:    {Path: "red-iso-half-cube.png"},
		ImgYellowCube2D:      {Path: "yellow-2d-cube.png"},
		ImgYellowCubeIso:     {Path: "yellow-iso-cube.png"},
		ImgYellowHalfCube2D:  {Path: "yellow-2d-half-cube.png"},
		ImgYellowHalfCubeIso: {Path: "yellow-iso-half-cube.png"},
	}

	for id, res := range imageResources {
		loader.ImageRegistry.Set(id, res)
		loader.LoadImage(id)
	}
}

func OpenAssetFunc(path string) io.ReadCloser {
	f, err := gameAssets.Open("resources/" + path)
	if err != nil {
		panic(err)
	}
	return f
}

//go:embed all:resources
var gameAssets embed.FS
