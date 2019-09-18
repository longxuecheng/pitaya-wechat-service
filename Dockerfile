FROM alpine:3.8

WORKDIR /geluxiya/gotrue/service

COPY ./gotrue ./

EXPOSE 8082

ENV ENV=prod

CMD ./gotrue
