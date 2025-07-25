# ------------------------
# Этап сборки
# ------------------------
FROM golang:1.24 AS build

# Объявляем аргумент ДО использования
ARG SERVICE_NAME
ARG LDFLAGS

# Определяем переменные окружения на основе ARG
ENV BIN_FILE="/opt/${SERVICE_NAME}/${SERVICE_NAME}-app"
ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

# Кэшируем зависимости
COPY go.mod .
COPY go.sum .
RUN go mod download

# Копируем исходный код
COPY . ${CODE_DIR}

# Сборка статического бинарника
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/server/*

# ------------------------
# Финальный образ
# ------------------------
FROM alpine:3.22.0

ARG SERVICE_NAME

# Метаданные
LABEL SERVICE="${SERVICE_NAME}"
LABEL MAINTAINERS="eandreev.dev@gmail.com"

# Копируем entrypoint
COPY ./build/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Переменные и бинарник
ENV BIN_FILE="/opt/${SERVICE_NAME}/${SERVICE_NAME}-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

# Конфиг
ENV CONFIG_FILE="/etc/${SERVICE_NAME}/config.toml"
COPY ./configs/config.toml ${CONFIG_FILE}

ENTRYPOINT ["/entrypoint.sh"]