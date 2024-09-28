# ToDo

- [ ] fix docker build, add tests from all repos to include all requirements
- [ ] eliminate go-echo-swagger, upload swagger files to apidocs.<domain>.com/service. Go-echo-swagger adds a lot of dependencies and slows down every request
- [ ] [server/middleware/ip_address.go#L18](use echo.RealIP()) instead of our custom code
- [ ] Set up proper middleware to get user ip address correctly, https://echo.labstack.com/docs/ip-address
- [ ] refactor middleware to have config same as [echo ones](https://github.com/labstack/echo/blob/master/middleware/body_dump.go#L18)


## Guidelines
### How to update ALL go.mod files?
` ./make.sh update-version <old-version> <new-version>`
For example
 `./make.sh update-version v0.0.0-20210101163630-b4ea9f10773c v0.0.0-20240928194423-e378b7eda3d5`