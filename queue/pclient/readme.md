## Test

### Set up emulator
- install emulator `gcloud components install pubsub-emulator`
- read env file `source .local.env`
- start emulator `gcloud beta emulators pubsub start --project=$PUBSUB_PROJECT_ID`
- set PUBSUB_EMULATOR_HOST env to the enumalor `$(gcloud beta emulators pubsub env-init)`
- dont forget to create topic using by using [CreateTopic](./client.go) client method

## ToDo
- [ ] pclient.Publish do NOT return messageID all the time - in some case it just an empty string
- [ ] Do not parse config file all the time use configurator
