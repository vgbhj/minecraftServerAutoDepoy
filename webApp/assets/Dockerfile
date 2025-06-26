FROM openjdk:21-jdk-slim

WORKDIR /data

COPY . /data

CMD ["java", "-Xmx1G", "-Xms1G", "-jar", "server.jar", "nogui"]