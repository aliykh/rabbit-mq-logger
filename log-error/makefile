CURRENT_DIR=$(shell pwd)

APP_NAME = app
APP_CMD_DIR=${CURRENT_DIR}/cmd
APP_BIN_DIR=${CURRENT_DIR}/bin


run:
	go build -o ${APP_BIN_DIR}/${APP_NAME} ${APP_CMD_DIR} && ${APP_BIN_DIR}/./${APP_NAME}