docker stop ecommerce && docker rm ecommerce &
jet -dsn="file://${PWD}/database.sqlite" -schema=main -path=./db &
templ generate &
wait
CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/go-ecommerce -a server.go
docker build . -t go-ecommerce:1
docker run -d -p 3000:3000 --name ecommerce go-ecommerce:1