deploy:
	GOOS=linux GOARCH=amd64 go build && \
		zip brexit-progress.zip brexit-progress -x Makefile && \
		aws lambda update-function-code --function-name brexit-progress --zip-file fileb://brexit-progress.zip && \
		rm brexit-progress.zip
