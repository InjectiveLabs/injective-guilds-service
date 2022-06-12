# Injective Trading Guilds service

An off-chain service for guilds-related queries

Install the injective-guilds binary

```
make install
```

Authorize the master from the default member through [MsgGrant](https://api.injective.exchange/#authz-msggrant) for the following messages:

```
"/injective.exchange.v1beta1.MsgCreateSpotLimitOrder",
"/injective.exchange.v1beta1.MsgCreateSpotMarketOrder",
"/injective.exchange.v1beta1.MsgCancelSpotOrder",
"/injective.exchange.v1beta1.MsgBatchUpdateOrders",
"/injective.exchange.v1beta1.MsgBatchCancelSpotOrders",
"/injective.exchange.v1beta1.MsgDeposit",
"/injective.exchange.v1beta1.MsgWithdraw",
"/injective.exchange.v1beta1.MsgCreateDerivativeLimitOrder",
"/injective.exchange.v1beta1.MsgCreateDerivativeMarketOrder",
"/injective.exchange.v1beta1.MsgCancelDerivativeOrder",
"/injective.exchange.v1beta1.MsgBatchUpdateOrders",
"/injective.exchange.v1beta1.MsgBatchCancelDerivativeOrders",
"/injective.exchange.v1beta1.MsgDeposit",
"/injective.exchange.v1beta1.MsgWithdraw"
```

Run this once to update the denom in the newly created db

```
# this will start mongo + init replset
make dev

# copy .env.example to .env and fill values
cp .env.example .env
```

To create a guild

```
injective-guilds add-guild \
--derivative-id=0x54d4505adef6a5cef26bc403a33d595620ded4e15b9e2bc3dd489b714813366a --derivative-require=1000 \
--capacity=150 --name "Hades Raiders" --description "Hades Raiders Guild" --master=inj14m8wrpeerjfjmutl7lzyvf48myx4lcrc75rtnl \
--default-member=inj14rhj922slkuczyzu7ah45pm84904ujdnjlnjcc --exchange-url=sentry0.injective.network:9910 \
--db-url=mongodb://mongo:27017 --lcd-url=https://lcd.injective.network
```

To delete a guild

```
injective-guilds delete-guild --guild-id=<HEX_STRING>
```

Start the api

```
injective-guilds api
```

Start the process (on another terminal)

```
injective-guilds process
```

## Deploy service on a cloud instance:

Initialize the services

```
mkdir -p ~/injective
cd injective && git clone https://github.com/InjectiveLabs/injective-guilds-service.git

cd injective-guilds-service
git checkout dev && git pull
rm -rf deployment/var
# build injective-guild binary
APP_ENV=test docker-compose -f deployment/testnet.yaml build injective-guilds-api
# setup mongo db
APP_ENV=test docker-compose -f deployment/testnet.yaml up -d mongo
APP_ENV=test docker-compose -f deployment/testnet.yaml up -d mongo-setup && sleep 10
# up guilds apps
APP_ENV=test docker-compose -f deployment/testnet.yaml up -d injective-guilds-api injective-guilds-process
```

Run commands

```
Example 1

docker exec -it injective-guilds-api injective-guilds add-guild \
--derivative-id 0x8158e603fb80c4e417696b0e98765b4ca89dcf886d3b9b2b90dc15bfb1aebd51 \
--derivative-require=20 \
--capacity=150 \
--name=Zeus --description "Akukx Guild" \
--master=inj13q8u96uftm0d7ljcf6hdp0uj5tyqrwftmxllaq \
--default-member=inj14au322k9munkmx5wrchz9q30juf5wjgz2cfqku \
--exchange-url=sentry1.injective.dev:9910 --db-url=mongodb://mongo:27017 \
--lcd-url=https://testnet.lcd.injective.dev

Example 2

docker exec -it injective-guilds-api injective-guilds add-guild \
--spot-id=0x0511ddc4e6586f3bfe1acb2dd905f8b8a82c97e1edaef654b12ca7e6031ca0fa --name "Hades" \
--description "Injective guild" --spot-require=10/20 --master=inj1wng2ucn0ak3aw5gq9j7m2z88m5aznwntqnekuv \
--default-member=inj1awx03zmnnlsjuvp7x8ac3lphw50p0nea6p2584 \
--exchange-url=sentry1.injective.dev:9910 --db-url=mongodb://mongo:27017 \
--lcd-url=https://testnet.lcd.injective.dev

# delete a guild
docker exec -it injective-guilds-api injective-guilds delete-guild --guild-id=<guild_id> --db-url=mongodb://mongo:27017
```
