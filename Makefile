gen-goa: export GOPROXY=direct
gen-goa:
	rm -rf ./api/gen
	go generate ./api/...
gen: gen-goa
dev-docker:
	docker run -d -p 27017:27017 --name mongo mongo
dev-docker-off:
	docker kill mongo redis 
	docker rm mongo redis
install:
	go install github.com/InjectiveLabs/injective-guilds-service/cmd/injective-guilds/...
dev:
	mkdir -p var/mongo/
	mongod --dbpath ./var/mongo > var/mongo/output.txt & echo $$! > var/mongo/mongod.pid
dev-off:
	kill -9 `cat ./var/mongo/mongod.pid`
