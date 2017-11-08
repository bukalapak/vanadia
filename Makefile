vanadia: main.go blueprint/drafter.go blueprint/* postman/* config/* vendor/*
	@go build -o vanadia

clean:
	rm -f vanadia
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
