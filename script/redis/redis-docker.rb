require "dotenv"

Dotenv.load(".dev.env")

$config = {
  servicename: ENV["REDIS_SERVICENAME"],
  host_port: ENV["REDIS_HOST_PORT"],
  docker_port: ENV["REDIS_DOCKER_PORT"],
  image_name: ENV["REDIS_IMAGE_NAME"],
  redis_db: ENV["REDIS_DB"],
  host_volume: ENV["REDIS_HOST_VOLUME"],
  docker_volume: ENV["REDIS_DOCKER_VOLUME"]
}

$command = {
  make: <<~CMD,
    docker run --name #{$config[:servicename]} --network infrastructure_network --restart=on-failure:10 \
    -p #{$config[:host_port]}:#{$config[:docker_port]} \
    -v #{$config[:host_volume]}:#{$config[:docker_volume]} \
    --sig-proxy=false -d #{$config[:image_name]}
  CMD
  start: "docker start #{$config[:servicename]}",
  stop: "docker stop #{$config[:servicename]}",
  debug: "iredis --rainbow",
  log: "docker logs -t #{$config[:servicename]}",
  prune: "docker rm -f #{$config[:servicename]}",
  stat: "docker stats #{$config[:servicename]}"
}

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

  system($command[command])
end

main
