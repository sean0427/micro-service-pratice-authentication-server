FROM golang:1.19-alpine as build

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0  go build -o /docker-exec ./cmd/main.go

#####################################

FROM alpine:3.16.3

ENV APP_USER app
ENV APP_HOME /go/src/app
 
RUN addgroup -S $APP_USER && adduser -S $APP_USER -G $APP_USER
RUN mkdir -p $APP_HOME


WORKDIR $APP_HOME

COPY --from=build /docker-exec $APP_HOME/docker-exec
RUN chown -R $APP_USER:$APP_USER $APP_HOME
RUN chmod -R 100 $APP_HOME
EXPOSE 8080
USER $APP_USER

CMD [ "./docker-exec" ]