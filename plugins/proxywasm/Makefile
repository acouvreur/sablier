build:
	go generate
	tinygo build -o sablierproxywasm.wasm -scheduler=none -target=wasi ./

docker:
	docker build -t acouvreur/sablier-proxy-wasm:latest .