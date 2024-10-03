FROM golang:1.22.2-alpine AS builder
WORKDIR /usr/local/src
RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["../go.mod", "../go.sum", "./"]
RUN go mod download

COPY ./ ./
RUN go build -o ./bin/go .

FROM alpine
COPY .env .env
COPY stock_dml.sql stock_dml.sql
COPY --from=builder /usr/local/src/bin/go /
RUN apk --no-cache add bash
EXPOSE 7557
CMD [ "/go" ]