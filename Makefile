gen-goa: export GOPROXY=direct
gen-goa:
	rm -rf ./api/gen
	go generate ./api/...
gen: gen-goa
install:
	go install github.com/InjectiveLabs/injective-guilds-service/cmd/injective-guilds/...
dev:
	mkdir -p var/mongo/
	mongod --replSet rs0 --dbpath ./var/mongo > var/mongo/output.txt & echo $$! > var/mongo/mongod.pid
	echo "Waiting 5s before initiating Replica Set.." && sleep 5;
	(mongo --eval "rs.status()" | grep "NotYetInitialized") && mongo --eval "rs.initiate()"
dev-off:
	kill -9 `cat ./var/mongo/mongod.pid`
