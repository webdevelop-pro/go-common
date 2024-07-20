# ToDo

- [ ] fix docker build, add tests from all repos to include all requirements
- [ ] eliminate go-echo-swagger, upload swagger files to apidocs.<domain>.com/service. Go-echo-swagger adds a lot of dependencies and slows down every request
- [ ] [server/middleware/ip_address.go#L18](use echo.RealIP()) instead of our custom code
- [ ] Set up proper middleware to get user ip address correctly, https://echo.labstack.com/docs/ip-address
- [ ] refactor middleware to have config same as [echo ones](https://github.com/labstack/echo/blob/master/middleware/body_dump.go#L18)
