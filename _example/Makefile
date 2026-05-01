PROJECT_ID     ?= $(shell gcloud config get-value project)
REGION         ?= asia-northeast1
SERVICE        ?= govite-example
KO_DOCKER_REPO ?= gcr.io/$(PROJECT_ID)

.PHONY: deploy build-frontend

## deploy: フロントエンドをビルドし、ko でイメージをビルド・プッシュして Cloud Run にデプロイする
deploy: build-frontend
	IMAGE=$$(KO_DOCKER_REPO=$(KO_DOCKER_REPO) ko build ./entrypoint/prod) && \
	gcloud run deploy $(SERVICE) \
		--image=$$IMAGE \
		--region=$(REGION) \
		--platform=managed \
		--allow-unauthenticated

## build-frontend: npm パッケージのインストール、コード生成、TypeScript コンパイル、Vite ビルドを行う
build-frontend:
	npm install
	go run entrypoint/dev/main.go -gen
	npx tsc -b
	npx vite build
