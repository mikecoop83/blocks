mkdir -p out/
rm -f out/*
env GOOS=js GOARCH=wasm go build -o out/blocks.wasm github.com/mikecoop83/blocks
cp resources/wasm_exec.js out/
cp resources/blocks.html out/
