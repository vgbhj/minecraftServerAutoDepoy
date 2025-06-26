FROM openjdk:21-jdk-slim

WORKDIR /data

COPY . /data

CMD ["java", "-Xmx512MB", "-Xms512MB", "-jar", "server.jar", "nogui"]
