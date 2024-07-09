# ToDo

- [ ] fix docker build, add tests from all repos to include all requirements
- [ ] server/middleware/log_request, provide a way to set up log level
- [ ] eliminate go-echo-swagger, upload swagger files to docs.<domain>.com/service. Go-echo-swagger adds a lot of dependencies and slows down every request
- [ ] [server/middleware/ip_address.go#L18](use echo.RealIP()) instead of our custom code
- [ ] Set up proper middleware to get user ip address correctly, https://echo.labstack.com/docs/ip-address
- [ ] move consts from keys/keys.go to the package where they used

