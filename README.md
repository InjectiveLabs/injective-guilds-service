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
	--derivative-id=0xc559df216747fc11540e638646c384ad977617d6d8f0ea5ffdfc18d52e58ab01 \
	--spot-id=0xfbc729e93b05b4c48916c1433c9f9c2ddb24605a73483303ea0f87a8886b52af \
	--name=testguild --description "a test guild" --master=inj1awx03zmnnlsjuvp7x8ac3lphw50p0nea6p2584 \
	--default-member=inj1zggdm44ln2gu7c5d2ge4wyr4wfs0cfn5lyfw4k --exchange-url=sentry2.injective.network:9910
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
