version := 0.0.2
projectName := uggo

build: format compile document push

format:
	go fmt ./...

compile:
	go build uggo.go

document:
	echo "# uggo" > README.md
	echo "helper utilities for uggly ecosystem" >> README.md
	echo "" >> README.md
	echo '```' >> README.md
	go doc -all uggo >> README.md
	echo '```' >> README.md

push:
	git add uggo.go
	git add README.md
	git add Makefile
	git commit -m "$(M)"
	git tag v$(version)
	git push origin v$(version)
	git push origin master
