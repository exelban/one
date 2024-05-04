VERSION = 0.0.0

.SILENT: release
.PHONY: release

release:
	rm -rf bin && rm -rf release && mkdir bin && mkdir release
	echo "Building release version $(VERSION):"

	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o bin/one && tar -czf release/one_$(VERSION)_darwin_x86_64.tar.gz -C bin one && rm bin/one
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o bin/one && tar -czf release/one_$(VERSION)_darwin_arm64.tar.gz -C bin one && rm bin/one

	echo "darwin completed."

	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o bin/one && tar -czf release/one_$(VERSION)_linux_x86_64.tar.gz -C bin one && rm bin/one
	GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$(VERSION)" -o bin/one && tar -czf release/one_$(VERSION)_linux_x86.tar.gz -C bin one && rm bin/one
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o bin/one && tar -czf release/one_$(VERSION)_linux_arm64.tar.gz -C bin one && rm bin/one

	echo "linux completed."

	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o bin/one && tar -czf release/one_$(VERSION)_windows_x86_64.tar.gz -C bin one && rm bin/one
	GOOS=windows GOARCH=386 go build -ldflags "-X main.version=$(VERSION)" -o bin/one && tar -czf release/one_$(VERSION)_windows_x86.tar.gz -C bin one && rm bin/one
	GOOS=windows GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o bin/one && tar -czf release/one_$(VERSION)_windows_arm64.tar.gz -C bin one && rm bin/one

	echo "windows completed."
	rm -rf bin
