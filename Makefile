build:
	go build -o email-notification
archbuild:
	docker build -t builder/package . && \
	docker create --name builder builder/package && \
	docker cp builder:/builds/output "$$(pwd)" && \
	docker rmi -f builder/package
