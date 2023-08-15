require "dotenv"
require "json"
require "mongo"

Dotenv.load(".dev.env")

$config = {
  servicename: ENV["MONGO_SERVICENAME"],
  host_port: ENV["MONGO_HOST_PORT"],
  docker_port: ENV["MONGO_DOCKER_PORT"],
  image_name: ENV["MONGO_IMAGE_NAME"],
  mongo_user: ENV["MONGO_INITDB_ROOT_USERNAME"],
  mongo_pass: ENV["MONGO_INITDB_ROOT_PASSWORD"],
  mongo_db: ENV["MONGO_INITDB_DATABASE"],
  host_volume: ENV["MONGO_HOST_VOLUME"],
  docker_volume: ENV["MONGO_DOCKER_VOLUME"],
  mongo_path: ENV["MONGO_PATH"]
}

$command = {
  make: <<~CMD,
    docker run  --network infrastructure_network --name #{$config[:servicename]} --restart=on-failure:10 \
    -e TITLE=#{$config[:servicename]} \
    -e MONGO_INITDB_ROOT_USERNAME=#{$config[:mongo_user]} \
    -e MONGO_INITDB_ROOT_PASSWORD=#{$config[:mongo_pass]} \
    -e MONGO_INITDB_DATABASE=#{$config[:mongo_db]} \
    -p #{$config[:host_port]}:#{$config[:docker_port]} \
    -v #{$config[:host_volume]}:#{$config[:docker_volume]} \
    --sig-proxy=false -d #{$config[:image_name]}
  CMD
  start: "docker start #{$config[:servicename]}",
  stop: "docker stop #{$config[:servicename]}",
  debug: "mongosh #{$config[:mongo_path]}",
  log: "docker logs -t #{$config[:servicename]}",
  prune: "docker rm -f #{$config[:servicename]}",
  stat: "docker stats #{$config[:servicename]}"
}

def initMongo
  client = Mongo::Client.new(ENV["MONGO_DB_URI"])

  db_name = ENV["MONGO_DB_NAME"]
  collection_name = "recipes"
  db = client.database
  collection = db[collection_name]
  file = File.read(ENV["MONGO_SAMPLE_DATA"])
  recipes = JSON.parse(file)
  result = collection.insert_many(recipes)
  puts("Inserted recipes: #{result.inserted_count}")
end

def main
  if ARGV.empty?
    puts("Please provide a command as the first argument.")
    return
  end

  command = ARGV[0].to_sym

  unless $command.key?(command)
    puts("Invalid command: #{command}")
    return
  end

  case command
  when :make
    system($command[command])
    sleep(1)
    initMongo
  else
    system($command[command])
  end
end

main
