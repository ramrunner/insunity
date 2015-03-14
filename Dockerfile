FROM ubuntu:14.04
MAINTAINER anon ymous, dsp@2f30.org
COPY bashrc /root/.bashrc
COPY insunity /root/go/src/insunity
ENV GOBIN /root/go/bin
ENV GOPATH /root/go
ENV PATH /root/go/bin:$PATH
ENV HOME /root
WORKDIR /root
RUN \
  apt-get update && \
  apt-get install -y golang git && \
  apt-get install -y protobuf-compiler && \
  mkdir -p /root/go/src && \
  mkdir -p /root/go/bin && \
  go get github.com/chai2010/protorpc && \
  go get github.com/chai2010/protorpc/protoc-gen-go

RUN protoc --go_out=go/src/insunity/entropy.pb/ go/src/insunity/entropy.pb/entropy.proto
RUN go install /root/go/src/insunity/insunity.go

CMD ["bash"]
