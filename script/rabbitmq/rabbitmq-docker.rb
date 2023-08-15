require "dotenv"

Dotenv.load(".dev.env")

$config = {
  servicename: ENV["RABBITMQ_SERVICENAME"],
  host_port_dashboard: ENV["RABBITMQ_HOST_PORT_DASHBOARD"],
  docker_port_dashboard: ENV["RABBITMQ_DOCKER_PORT_DASHBOARD"],
  host_port_service: ENV["RABBITMQ_HOST_PORT_SERVICE"],
  docker_port_service: ENV["RABBITMQ_DOCKER_PORT_SERVICE"],
  image_name: ENV["RABBITMQ_IMAGE_NAME"],
  rabbit_user: ENV["RABBITMQ_DEFAULT_USER"],
  rabbit_pass: ENV["RABBITMQ_DEFAULT_PASS"]
}

$command = {
  make: <<~CMD,
    docker run  --network infrastructure_network --name  #{$config[:servicename]} --restart=on-failure:10 \
    -e RABBITMQ_DEFAULT_USER=#{$config[:rabbit_user]} \
    -e RABBITMQ_DEFAULT_PASS=#{$config[:rabbit_pass]} \
    -p #{$config[:host_port_dashboard]}:#{$config[:docker_port_dashboard]} \
    -p #{$config[:host_port_service]}:#{$config[:docker_port_service]} \
    --sig-proxy=false -d #{$config[:image_name]}
  CMD
  start: "docker start #{$config[:servicename]}",
  stop: "docker stop #{$config[:servicename]}",
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

  puts($command[command])
  system($command[command])
end

main
