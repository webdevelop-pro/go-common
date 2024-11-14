# ToDo
- [ ] improve healtcheck checker:
- - [ ] based on this article https://testfully.io/blog/api-health-check-monitoring/
- - [ ] add (rename) endpoing `/health` should return generic health status 200 ok, memory consumption, dependencies statuses (db connection, db connection amount, pubsub connection etc), 3rd party dependencies healthcheck status
- - [ ] add endpoint `/health/live` should return 200 when service is ready to accept connections
- - [ ] add endpoint `/health/ready` when service is up but connecting to 3rd party dependencies
- [ ] replace swagger library (currently pollute go.mod file and inspect every request)

