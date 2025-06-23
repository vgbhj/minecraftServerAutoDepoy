FROM openjdk:21-jdk-slim

WORKDIR /data

COPY . /data

CMD ["java", "-Xmx2G", "-Xms2G", "-jar", "server.jar", "nogui"]
