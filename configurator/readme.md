# Configuration
Read data from disk using ENV_FILE env variable.
This operation is `slow`, use it during process start to load envs from file

```golang
  var cfg Config
  err := configurator.NewConfiguration(&cfg, "")
	if err != nil {
		panic(err)
	}
```

# Configurator
Read data from the configuration, do not try to read ENV_FILE
Use this method in running process to provide envs for services

```golang
	conf := configurator.NewConfigurator()
	cfg := conf.New("logger", &Config{}).(*Config)
```


# ToDo
- [ ] create a singleton so we can read file from disk once and use cached data after
- [ ] find out a way not to use github.com/jinzhu/copier. Do we really need [this code?](./configurator.go#78)