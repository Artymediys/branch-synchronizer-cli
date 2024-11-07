#!/bin/bash

# Название приложения и директория сборки
APP_NAME="bsync_cli"
APP_DIR="./bin"

# Функция для сборки под конкретную ОС и архитектуру
build() {
  OS=$1
  ARCH=$2
  OUTPUT_DIR=$3
  OUTPUT_FILE=$4

  echo "Building for $OS $ARCH..."
  GOOS=$OS GOARCH=$ARCH go build -o "$OUTPUT_DIR"/"$OUTPUT_FILE"
  if [ $? -ne 0 ]; then
    echo "Build failed for $OS $ARCH"
    exit 1
  fi
}

# Сборка в зависимости от аргументов
if [ "$1" == "all" ]; then
  build "windows" "amd64" $APP_DIR/windows "${APP_NAME}_amd64.exe"
  build "linux" "amd64" $APP_DIR/linux "${APP_NAME}_amd64"
  build "darwin" "amd64" $APP_DIR/macos "${APP_NAME}_amd64"
  build "darwin" "arm64" $APP_DIR/macos "${APP_NAME}_arm64"
elif [ "$1" == "windows" ]; then
  build "windows" "amd64" $APP_DIR/windows "${APP_NAME}_amd64.exe"
elif [ "$1" == "linux" ]; then
  build "linux" "amd64" $APP_DIR/linux "${APP_NAME}_amd64"
elif [ "$1" == "macos" ]; then
  build "darwin" "amd64" $APP_DIR/macos "${APP_NAME}_amd64"
  build "darwin" "arm64" $APP_DIR/macos "${APP_NAME}_arm64"
else
  # Сборка под текущую ОС и архитектуру
  go build -o "$APP_DIR"/"$APP_NAME"
fi