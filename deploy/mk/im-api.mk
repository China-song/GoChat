VERSION=latest

SERVICE_NAME=im
SERVICE_TYPE=api

# 测试环境配置
# docker的镜像发布地址
DOCKER_REPO_TEST=registry.cn-hangzhou.aliyuncs.com/godocker-song/${SERVICE_NAME}-${SERVICE_TYPE}-dev
# 测试版本
VERSION_TEST=$(VERSION)
# 编译的程序名称
APP_NAME_TEST=go-chat-${SERVICE_NAME}-${SERVICE_TYPE}-test

# 测试下的编译文件
DOCKER_FILE_TEST=./deploy/dockerfile/Dockerfile_${SERVICE_NAME}_${SERVICE_TYPE}_dev

# 测试环境的编译发布
build-test:

	set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&& go build -o bin/${SERVICE_NAME}-${SERVICE_TYPE} ./apps/${SERVICE_NAME}/${SERVICE_TYPE}/${SERVICE_NAME}.go
	docker build . -f ${DOCKER_FILE_TEST} --no-cache -t ${APP_NAME_TEST}

# 镜像的测试标签
tag-test:

	@echo 'create tag ${VERSION_TEST}'
	docker tag ${APP_NAME_TEST} ${DOCKER_REPO_TEST}:${VERSION_TEST}

publish-test:

	@echo 'publish ${VERSION_TEST} to ${DOCKER_REPO_TEST}'
	docker push $(DOCKER_REPO_TEST):${VERSION_TEST}

release-test: build-test tag-test publish-test