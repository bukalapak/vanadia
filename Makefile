apib-to-postman: main.go blueprint/* postman/*
	@go build -o apib-to-postman

clean:
	rm -f apib-to-postman
	rm -f blueprint/drafter.go
	rm -rf blueprint/ext

blueprint/drafter.go: blueprint/ext/drafter/bin/drafter
	@go get github.com/mjibson/esc
	go generate blueprint/parser.go

blueprint/ext/drafter/bin/drafter:
	mkdir -p blueprint/ext
	wget https://github.com/apiaryio/drafter/releases/download/v3.2.7/drafter-v3.2.7.tar.gz
	tar -xzf drafter-v3.2.7.tar.gz -C blueprint/ext
	rm drafter-v3.2.7.tar.gz
	mv blueprint/ext/drafter-v3.2.7 blueprint/ext/drafter
	$(MAKE) -C blueprint/ext/drafter drafter

.PHONY: clean
