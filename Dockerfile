FROM --platform=$BUILDPLATFORM golang:1.23.0 AS builder
ARG TARGETARCH
WORKDIR /workdir
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -trimpath -o switch-scraper .

FROM scratch
LABEL org.opencontainers.image.source=https://github.com/transacid/switch-scraper
LABEL org.opencontainers.image.description="a prometheus scraper that gets metrics from a netgear switch"
LABEL org.opencontainers.image.licenses=GPL-3.0-only
COPY --from=builder /workdir/switch-scraper ./
EXPOSE 8080
ENTRYPOINT [ "/switch-scraper" ]