FROM alpine

WORKDIR /app

COPY loggerApp .

CMD [ "/app/loggerApp" ]
