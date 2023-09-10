FROM debian:stable-20230904 as build

ARG LAVALINK_VERSION=3.7.8
ARG OPENJDK_VERSION=17

RUN mkdir /usr/app
WORKDIR /usr/app

# Installp depedencies
RUN apt-get update && apt-get install -y wget

# Download lavalink
RUN wget https://github.com/lavalink-devs/Lavalink/releases/download/${LAVALINK_VERSION}/Lavalink.jar

# Donwload open-jdk
RUN wget https://download.java.net/java/GA/jdk17/0d483333a00540d886896bac774ff48b/35/GPL/openjdk-${OPENJDK_VERSION}_linux-x64_bin.tar.gz && \
    tar xvf openjdk-${OPENJDK_VERSION}*_bin.tar.gz -C ./

FROM debian:stable-20230904-slim as prod

RUN mkdir /usr/app
WORKDIR /usr/app

# Edit PATH
ENV PATH="/usr/bin/jdk-17/bin:${PATH}"

# Copy binaries
COPY --from=build /usr/app/jdk-17 /usr/bin/jdk-17
COPY --from=build /usr/app/Lavalink.jar /usr/app

# Copy config file
COPY application.yml .

ENTRYPOINT ["java", "-jar", "Lavalink.jar"]
