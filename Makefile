user-rpc-dev:
	@make -f deploy/mk/user-rpc.mk release-test

user-api-dev:
	@make -f deploy/mk/user-api.mk release-test

release-test: user-rpc-dev user-api-dev

install-service:
	cd ./deploy/script && chmod +x ./release-test.sh && ./release-test.sh