# go-game-block-placement

Demo combining Ebitengine, EbitenUI, Resolv, ebitengine-resource to create a demo which allows you to place blocks in both a 2D and Isometric space.

![](docs/block-placement.gif?raw=true)

Build WASM:
```shell
env GOOS=js GOARCH=wasm go build -o blockPlacement-1-0-1.wasm github.com/timothy-ch-cheung/go-game-block-placement
```