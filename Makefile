.PHONY: build-recipes unit-test-recipes unit-test-rss tc-recipes tcr-recipes swag-recipes servedoc-recipes \
        build-rss-ph tc-rss-p tcr-rss-p swag-rss-p servedoc-rss-p \
        mongo-create mongo-start mongo-stop mongo-prune mongo-stat mongo-log mongo-debug mongo-init \
        redis-make redis-start redis-stop redis-debug redis-prune \
        vault-make vault-start vault-stop vault-prune vault-stat vault-log \
        init help run-recipes run-rss-p


# [Main] Build and Run 
init: ## Initialize this project by install tools and other package associate with this package infrastructure 
	@cd /path/to/directory/containing/files && bash ./script/setup.sh

help: ## Generate help
	@bash ./script/print_help.sh

run-recipes: build-recipes ## Build and run recipes services
	@cd Recipes_Services && ./bin/main

run-rss-p: build-rss-p ## Build and run Rss services producer
	@cd Rss_Producer_Services && ./bin/main

run-rss-c: build-rss-c ## Build and run Rss  services worker
	@cd Rss_Worker_Services && ./bin/main

build-all-docker: build-recipes-image build-rss-p build-rss-c  ## Build docker images  for all services  

run-project: init build-all-docker  ## Install project dependency , build and run the project
	docker compose up


# [Secure Check]
secure-scan-recipes:  ## Run snyk test on  recipes services 
	@cd Recipes_Services && snyk test

secure-scan-rss-p: ## Run snyk test on rss producer service 
	@cd Rss_Producer_Services &&  snyk test 

secure-scan-rss-c: ## Run snyk test on rss consumer service 
	@cd Rss_Worker_Services && snyk test

secure-scan-all: secure-scan-recipes secure-scan-rss-c secure-scan-rss-p


# [Recipes Services] specific 
build-recipes: ## Build project into bin directory for recipes service 
	@cd Recipes_Services && go build -o bin/main cmd/server/main.go

unit-test-recipes: ## Run all the unit test for recipes service 
	@cd Recipes_Services && go test -v -run ".*Unit" ./...

tc-recipes: ## Test success and commit  recipes service
	@cd Recipes_Services &&  go  test ./... & git commit -am "Work in progress 📝"  --amend  --no-edit

tcr-recipes: ## Test success and commit or reverse recipes service 
	@cd Recipes_Services && go test ./... && git commit -am "Work in Progress 📝" || git reset --hard

swag-recipes: ## Generate swag file 
	@cd Recipes_Services && swag init -g cmd/server/main.go --outputTypes json

servedoc-recipes: swag ## Generate swag file and serve doc in one go  for recipes service 
	@cd Recipes_Services && swagger serve ./docs/swagger.json

build-recipes-image: ## Build Recipes services docker image
	@cd Recipes_Services && docker image build --no-cache --tag recipes_app .


# [Rss-producer Services] Specific
unit-test-rss-p: ## Run all the unit test for rss service  
	@cd Rss_Producer_Services && go test -v run ".*Unit" ./...

build-rss-p: ## Build project into bin directory for rss service  prodpucer 
	@cd Rss_Producer_Services && go build -o bin/main cmd/server/main.go

tc-rss-p: ## Test success and commit for rss service producer
	@cd Rss_Producer_Services &&  go  test ./... & git commit -am "Work in progress 📝"  --amend  --no-edit

tcr-rss-p: ## Test success and commit or reverse for rss service producer
	@cd Rss_Producer_Services && go test ./... && git commit -am "Work in Progress 📝" || git reset --hard

swag-rss-p: ## Generate swag file  for rss service producer
	@cd Rss_Producer_Services && swag init -g cmd/server/main.go --outputTypes json

servedoc-rss-p: swag ## Generate swag file and serve doc in one go  for rss service producer
	@cd Rss_Producer_Services && swagger serve ./docs/swagger.json

build-rss-p-images: ## Build Rss Producer Service docker image
	@cd Rss_Producer_Services && docker image build --no-cache --tag rss_producer .


# [Rss-Consumer Services] Specific 
build-rss-c: ## Build project into bin directory for rss service worker
	@cd Rss_Worker_Services && go build -o bin/main cmd/server/main.go

tc-rss-c: ## Test success and commit for rss service worker
	@cd Rss_Worker_Services &&  go  test ./... & git commit -am "Work in progress 📝"  --amend  --no-edit

tcr-rss-c: ## Test success and commit or reverse for rss service worker
	@cd Rss_Worker_Services && go test ./... && git commit -am "Work in Progress 📝" || git reset --hard

swag-rss-c: ## Generate swag file  for rss service worker
	@cd Rss_Worker_Services && swag init -g cmd/server/main.go --outputTypes json

servedoc-rss-c: swag ## Generate swag file and serve doc in one go  for rss service worker
	@cd Rss_Worker_Services && swagger serve ./docs/swagger.json

build-rss-c-image: ## Build RSS consumer service  docker image 
	@cd Rss_Worker_Services && docker image build --tag rss_worker:latest . 


# [Docker] Mongo
mongo-create: ## Create mongo docker container 
	@ruby ./script/mongo/mongo-docker.rb  make

mongo-start: ## Start stopped mongo docker container 
	@ruby ./script/mongo/mongo-docker.rb  start

mongo-stop: ## Stop mongo docker container 
	@ruby ./script/mongo/mongo-docker.rb stop

mongo-prune: ## Remove mongo docker container and its volume 
	@ruby ./script/mongo/mongo-docker.rb prune

mongo-stat:  ## Printout stat info of mongo docker container 
	@ruby ./script/mongo/mongo-docker.rb stat

mongo-log: ## Printout log of mongo docker container 
	@ruby ./script/mongo/mongo-docker.rb log

mongo-debug: ## exec mongosh  to debug mongo docker container 
	@ruby ./script/mongo/mongo-docker.rb debug

mongo-init: ## initialize data into mongo docker container 
	@ruby ./script/initmongo.rb


# [Docker] Redis
redis-make: ## Make redis docker  container 
	@ruby ./script/redis/redis-docker.rb make

redis-start: ## Start stopped redis docker  container 
	@ruby ./script/redis/redis-docker.rb start

redis-stop: ## Stopped redis docker  container 
	@ruby ./script/redis/redis-docker.rb stop

redis-debug:  ## Debug redis docker  container with iredis 
	@ruby ./script/redis/redis-docker.rb debug

redis-prune: ## remove redis docker  container and its volume
	@ruby ./script/redis/redis-docker.rb prune

# [Docker] Vault
vault-make: ## Make vault docker container 
	@ruby ./script/vault/vault-docker.rb make

vault-start: ## Start stopped vault docker container
	@ruby ./script/vault/vault-docker.rb start

vault-stop: ## Stop vault docker container 
	@ruby ./script/vault/vault-docker.rb stop 

vault-prune: ## Remove vault docker container and its volume 
	@ruby ./script/vault/vault-docker.rb  prune

vault-stat: ## Vault container status 
	@ruby ./script/vault/vault-docker.rb  stat

vault-log:  ## Printout log of vault docker container 
	@ruby ./script/vault/vault-docker.rb log

# [Docker] RabbitMq
rabbitmq-make: ## Make Rabbitmq  
	@ruby ./script/rabbitmq/rabbitmq-docker.rb  make

rabbitmq-start: ## Start Rabbitmq  
	@ruby ./script/rabbitmq/rabbitmq-docker.rb  start

rabbitmq-stop: ## Stop Rabbitmq 
	@ruby ./script/rabbitmq/rabbitmq-docker.rb stop 

rabbitmq-prune: ## Remove Rabbitmq and its associated file
	@ruby ./script/rabbitmq/rabbitmq-docker.rb prune  

rabbitmq-log: ## Print out log of rabbitmq 
	@ruby ./script/rabbitmq/rabbitmq-docker.rb log

rabbitmq-stat: ## Print out stat of rabbitmq 
	@ruby ./script/rabbitmq/rabbitmq-docker.rb stat 




