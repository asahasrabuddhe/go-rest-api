FROM golang as builder
WORKDIR /app
COPY . .

WORKDIR /app/Databases
RUN ls
RUN go build main.go

FROM ubuntu
# MAINTAINER Ayush

COPY  --from=builder /app/Databases/main  /goapps/main
RUN chmod +x /goapps/main
# ENV PORT 8080
EXPOSE 8080
ENTRYPOINT /goapps/main