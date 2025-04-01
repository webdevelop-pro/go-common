# ToDo
- [ ] add partial %any% for example "access_token": "plaid-access-token-%any%'
- [ ] add ability to check response only using assertify library
- [ ] Create setUp() and tearDown() methods that will allow you to define instructions that will be executed before and after each test method. 
- [ ] Add stack trace with a start a function which trigger test. As for example when error code does not match we have this response
```golang
    actions.go:36: 
      Error Trace:    /home/adams/projects/webdevelop-pro/go-common/tests/actions.go:36
                        /home/adams/projects/webdevelop-pro/go-common/tests/general.go:54
      Error:          Not equal: 
                      expected: 200
                      actual  : 404
```
    Probably I would like to see something like this:
```golang
    actions.go:36: 
      Error Trace:    /home/adams/projects/webdevelop-pro/projectName/tests/my_test.go:123 <---- start of the test
                        /home/adams/projects/webdevelop-pro/go-common/tests/actions.go:36
                          /home/adams/projects/webdevelop-pro/go-common/tests/general.go:54
      Error:          Not equal: 
                      expected: 200
                      actual  : 404
```
