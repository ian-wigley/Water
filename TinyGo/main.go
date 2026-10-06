package main

func main() {
	Initialise()
}

// https://www.fermyon.com/blog/optimizing-tinygo-wasm

// env GOOS=js GOARCH=wasm go build -o water.wasm
// compress with ....
//

// tinygo build -o ./web/main.wasm -target=wasm .

// tinygo build -target=wasm -opt=z -no-debug -o ./web/main.wasm .
// wasm-opt -Oz ./web/main.wasm -o ./web/main.opt.wasm
// brotli -9 ./web/main.opt.wasm
