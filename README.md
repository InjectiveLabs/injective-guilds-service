# Injective Trading Guilds service

An off-chain service for guilds-related queries

To debug locally, first can install if we haven't get injective-guilds binary

```
make install
```

Run this once to update denom in newly created db

```
make dev # this will start mongo + init replset
injective-guilds update-denom 
```

To manually create a guild:

```
injective-guilds add-guild \
--spot-id=0xf04d1b7acf40b331d239fcff7950f98a4f2ab7adb2ceb8f65aa32ac29455d7b4 --spot-require=0/20 \
--name=testguild --description "Peiyun guild" --master=inj1wng2ucn0ak3aw5gq9j7m2z88m5aznwntqnekuv \
--default-member=inj1fpmlw98jka5dc9cjrwurvutz87n87y45skvqkv \
--exchange-url=sentry2.injective.network:9910 --db-url=mongodb://localhost:27017
```

To delete a guild

```
injective-guilds delete-guild --guild-id=<HEX_STRING>
```

To start api/process:

```
injective-guilds api
```

On another terminal 

```
injective-guilds process
```
