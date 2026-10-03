package main

func main() {
	InitS()
}


// https://www.fermyon.com/blog/optimizing-tinygo-wasm

// env GOOS=js GOARCH=wasm go build -o water.wasm
// compress with ....
//


// tinygo build -o main.wasm -target=wasm .


// tinygo build -target=wasm -opt=z -no-debug -o main.wasm .
// wasm-opt -Oz main.wasm -o main.opt.wasm
// brotli -9 main.opt.wasm
