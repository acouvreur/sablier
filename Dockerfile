FROM golang:1.15.3-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o ondemand-service .
FROM scratch AS bin
COPY --from=build /src/ondemand-service /
ENTRYPOINT [ "/ondemand-service" ]