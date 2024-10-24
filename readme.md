# ToDo

- [ ] update `docker/` image and go-common image
- [ ] fix docker build, add tests from all repos to include all requirements
- [ ] eliminate go-echo-swagger, upload swagger files to apidocs.<domain>.com/service. Go-echo-swagger adds a lot of dependencies and slows down every request
- [ ] [server/middleware/ip_address.go#L18](use echo.RealIP()) instead of our custom code
- [ ] Set up proper middleware to get user ip address correctly, https://echo.labstack.com/docs/ip-address
- [ ] refactor middleware to have config same as [echo ones](https://github.com/labstack/echo/blob/master/middleware/body_dump.go#L18)
- [ ] replace `configurator.NewConfiguration` with `configurator.NewConfigurator` and run `configurator.NewConfiguration` only at start


## Guidelines
### How to update ALL go.mod files?

1. add replace to go.mod file to work locally, example `replace github.com/webdevelop-pro/go-common/queue => /home/go-common/queue`
2. once you finish run tests to verify go-common is working properly
3. if tests are passed commit and push
4. one pushed open queue (or any other package) and update requirements with latest push go get github.com/webdevelop-pro/go-common/context@<hash>
5. look on the new version hash in queue/go.mod and update dependencies everywhere ` ./make.sh update-version <old-version> <new-version>`
6. For example `./make.sh update-version v0.0.0-20210101163630-b4ea9f10773c v0.0.0-20240928194423-e378b7eda3d5`
