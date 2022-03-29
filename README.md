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

(There are too many params for now)

```
injective-guilds add-guild \
--derivative-id=0x8158e603fb80c4e417696b0e98765b4ca89dcf886d3b9b2b90dc15bfb1aebd51 --derivative-require=20 \ --derivative-id=0x1c79dac019f73e4060494ab1b4fcba734350656d6fc4d474f6a238c13c6f9ced --derivative-require=10 \
--name=Akukx --description "Akukx Thomas Guild" --master=inj1wng2ucn0ak3aw5gq9j7m2z88m5aznwntqnekuv \
--default-member=inj1kgpvzl2sjd527a7u5jj99j9pdple5050yavsd4 --exchange-url=sentry2.injective.network:9910 --db-url=mongodb://mongo:27017 --lcd-url=https://lcd.injective.network
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
