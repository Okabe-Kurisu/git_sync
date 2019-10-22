FROM ubuntu:18.04

ENV GIT_SSL_NO_VERIFY=true
ENV TZ=Asia/Shanghai

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*