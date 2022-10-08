FROM alpine

WORKDIR /app

COPY brokerApp .

CMD [ "/app/brokerApp" ]
