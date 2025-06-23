FROM openjdk:17-jdk-slim

WORKDIR /data

COPY . /data

RUN echo "eula=true" > eula.txt

CMD ["java", "-Xmx2G", "-Xms2G", "-jar", "server.jar", "nogui"]
