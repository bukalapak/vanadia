apib-to-postman:
	@go build -o apib-to-postman

clean:
	rm apib-to-postman

.PHONY: clean
