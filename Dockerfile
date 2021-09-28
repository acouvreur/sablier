FROM golang:1.17-alpine AS build

ENV PORT 10000

COPY . /go/src/ondemand-service
WORKDIR /go/src/ondemand-service

RUN go build -o /go/bin/ondemand-service

FROM alpine
EXPOSE 10000
COPY --from=build /go/bin/ondemand-service /go/bin/ondemand-service

ENTRYPOINT [ "/go/bin/ondemand-service" ]
CMD [ "--swarmMode=true" ]