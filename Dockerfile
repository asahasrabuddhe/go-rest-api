FROM ubuntu
# MAINTAINER Ayush
#RUN mkdir app
#RUN cd app/
#WORKDIR app
COPY  ./Databases/main /app/main
RUN chmod +x /app/main
# ENV PORT 8080
EXPOSE 8080
ENTRYPOINT /app/main