FROM ubuntu
MAINTAINER Ayush
#RUN mkdir app
#RUN cd app/
#WORKDIR app
COPY  ./go-rest-api /app/go-rest-api
RUN chmod +x /app/go-rest-api

ENV PORT 8080
EXPOSE 8080
#CMD ["pwd && ls -la cloud_native"]
ENTRYPOINT /app/go-rest-api