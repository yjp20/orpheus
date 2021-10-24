build:
	go build ./cmd/orpheus
	cd ui && npm run build
