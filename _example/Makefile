PROJECT_ID     ?= $(shell gcloud config get-value project)
REGION         ?= asia-northeast1
SERVICE        ?= govite-example
KO_DOCKER_REPO ?= gcr.io/$(PROJECT_ID)

.PHONY: deploy

## deploy: フロントエンドをビルドし、ko でイメージをビルド・プッシュして Cloud Run にデプロイする
deploy:
	npm install
	npm run build:frontend
	IMAGE=$$(KO_DOCKER_REPO=$(KO_DOCKER_REPO) ko build ./entrypoint/prod) && \
	gcloud run deploy $(SERVICE) \
		--image=$$IMAGE \
		--region=$(REGION) \
		--platform=managed \
		--allow-unauthenticated
