FROM golang:alpine
ARG VERSION=none
WORKDIR /app
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -a -installsuffix cgo -o main .

FROM scratch
WORKDIR /app
COPY --from=0 /app .
CMD ["./main"]