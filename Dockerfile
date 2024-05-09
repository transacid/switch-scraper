FROM --platform=$BUILDPLATFORM golang:1.22.3 as builder
ARG TARGETARCH
WORKDIR /workdir
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -trimpath -o switch-scraper .

FROM scratch
COPY --from=builder /workdir/switch-scraper ./
EXPOSE 8080
ENTRYPOINT [ "/switch-scraper" ]