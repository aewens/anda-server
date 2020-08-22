# Contributing

1. Fork the repository.
1. Create your feature branch using: `git checkout -b your-new-feature`
1. Add the appropriate tests to the respective `*_test.go` of the changed files.
1. Format the code using: `gofmt -w .`
1. Verify there are no linting warnings/errors using: `golangci-lint run`
1. Verify all tests pass: `CONFIG_PATH=/path/to/etc/config.json go test ./...`
1. Add your changed files: `git add .`
1. Commit your changes using: `git commit -m "DESCRIPTION"`
1. Push your changes to the feature branch: `git push origin your-new-feature`
1. Submit a pull request

Please make the pull request title and commit message(s) both descriptive and informative.
