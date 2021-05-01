.PHONY: env env-test env-ci clear-env

_PROJECT_NAME = bookmark

ifneq ("$(wildcard ${_PROJECT_DIR}/.env)", "")
	_ENV_FILE=.env
else
	_ENV_FILE=.env.dist
endif

env:
	@grep -h -s -v '^#' ${_ENV_FILE} | sed 's/^/export /' | sed 's/=/="/' | sed 's/$$/"/'
	@echo '# run this command to configure your shell:'
	@echo '# eval "$$(make env)"'

clear-env:
	@grep -h -s -v '^#' .env.dist .env .env.ci .env.test | sed 's/\=.*//g' | sort | uniq | sed 's/^/unset /'

run-server:
	@go run cmd/${_PROJECT_NAME}/main.go server
