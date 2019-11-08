FROM 25.2.20.4/rancher/rancher-server-ubuntu

ENV GIT_SSL_NO_VERIFY=true
ENV TZ=Asia/Shanghai

ENV GOPATH=/go PATH=/go/bin:/usr/local/go/bin:${PATH} SHELL=/bin/bash
RUN curl -sLfk https://nexus.ebcpaas.com/repository/pipeline-depend/golang/go1.12.10.linux-amd64.tar.gz | tar xzf - -C /usr/local

RUN mkdir -p /go/src/git-sync
WORKDIR /go/src/git-sync/
ADD . ./

RUN go build -o /bin/git-sync

CMD ["git-sync"]