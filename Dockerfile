FROM golang:1.22.2 as builder
WORKDIR /workdir
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY *.go ./
RUN CGO_ENABLED=0 go build -v -trimpath -o switch-monitor .

FROM scratch
COPY --from=builder /workdir/switch-monitor ./
EXPOSE 8080
ENTRYPOINT [ "/switch-monitor" ]