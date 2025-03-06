.PHONY: build clean deploy

# find all go files in src/internal/handlers/
GO_FILES := $(shell find src/internal/handlers -type f -name main.go)
# all go files after remove prefix "src/internal/handlers/" and subfix "/main.go"
CLEAN_GO_FILES := $(patsubst src/internal/handlers/%/main.go,%,$(GO_FILES))

build:
	go mod download
	$(foreach name, $(CLEAN_GO_FILES), \
		env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/$(name)/bootstrap ./src/internal/handlers/$(name)/main.go;\
		zip -j bin/$(subst /,_,$(name)).zip bin/$(name)/bootstrap;\
	)

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
