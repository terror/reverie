set dotenv-load

export EDITOR := 'nvim'

files := 'src/*'

default:
  just --list

ci: test lint forbid fmt-check

[group: 'dev']
dev-deps:
	brew install golangci-lint

[group: 'check']
forbid:
  ./bin/forbid

[group: 'format']
fmt:
  gofmt -w {{ files }}
  ./bin/retab

[group: 'check']
fmt-check:
	gofmt -l .

[group: 'check']
lint:
  golangci-lint run {{ files }}

[group: 'dev']
[script]
run:
	go run `fd .go ./src -E *_test.go`

[group: 'test']
test:
  go test ./src
