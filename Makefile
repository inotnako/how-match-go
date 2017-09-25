
byfile:
	echo '/etc/passwd\n/etc/hosts' | go run main.go -type file
	echo './testdata/file_with_go.txt\n./testdata/just_file.txt' | go run main.go -type file

byurl:
	echo 'https://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org\nhttps://golang.org' | go run main.go -type url
	echo 'https://golangsss.org\nhttps://golangsss.org\nhttps://golangsss.org\nhttps://golangxxx.org\nhttps://golangaaa.org' | go run main.go -type url
	echo 'https://github.com/' | go run main.go -type url

test:
	go test ./...

stress:
	go test -count=10 ./...
