FROM 25.2.20.4/rancher/rancher-server-ubuntu

ENV GIT_SSL_NO_VERIFY=true
ENV TZ=Asia/Shanghai


RUN apt-get update && apt-get install -y git wget gzip && rm -rf /var/lib/apt/lists/*
ENV GOPATH=/go PATH=/go/bin:/usr/local/go/bin:${PATH} SHELL=/bin/bash
RUN wget --no-check-certificate -O - https://nexus.ebcpaas.com/repository/pipeline-depend/golang/go1.12.10.linux-amd64.tar.gz | tar -xzf - -C /usr/local

RUN mkdir -p /go/src/git-sync
WORKDIR /go/src/git-sync/
ADD . ./

RUN go build -o /bin/git-sync

CMD ["git-sync"]