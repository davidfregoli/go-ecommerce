FROM alpine
WORKDIR /app
COPY ./bin/go-ecommerce ./
COPY ./database.sqlite ./
COPY ./assets/ ./assets/
CMD ["/app/go-ecommerce"]