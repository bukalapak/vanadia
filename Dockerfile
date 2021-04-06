FROM golang:1.12
WORKDIR /data
COPY . /go/src/github.com/SharperShape/vanadia
RUN cd /go/src/github.com/SharperShape/vanadia && bin/setup.sh
ENTRYPOINT [ "/go/src/github.com/SharperShape/vanadia/vanadia" ]
