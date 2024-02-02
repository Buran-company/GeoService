FROM golang:1.21

#FROM base as dev

#RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /home/hexedchild1/Kata/Repository/go-kata/course4Geoservice_1/

RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN go mod tidy
#RUN go get  -t -v ./...

# Build the application
#RUN go build -o main .

# Export necessary port
#EXPOSE 8080
#EXPOSE 5432
#CMD ["air"]