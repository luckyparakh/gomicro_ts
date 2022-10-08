FROM alpine

WORKDIR /app

COPY authApp .

CMD [ "/app/authApp" ]
