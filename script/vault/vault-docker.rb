require "dotenv"
require "erb"

Dotenv.load(".dev.env")

$config = {
  servicename: ENV["VAULT_SERVICENAME"],
  host_port: ENV["VAULT_HOST_PORT"],
  docker_port: ENV["VAULT_DOCKER_PORT"],
  image_name: ENV["VAULT_IMAGE_NAME"],
  vault_listen_address: ENV["VAULT_LISTEN_ADDRESS"],
  vault_api_address: ENV["VAULT_API_ADDRESS"],
  vault_storage_file_name: ENV["VAULT_STORAGE_FILE_NAME"],
  vault_host_volume: ENV["VAULT_HOST_VOLUME"],
  aws_kms_key_id: ENV["AWS_KMS_KEY_ID"],
  aws_region: ENV["AWS_REGION"],
  aws_access_key: ENV["AWS_ACCESS_KEY"],
  aws_secret_key: ENV["AWS_SECRET_KEY"],
  mongo_username: ENV["MONGO_INITDB_ROOT_USERNAME"],
  mongo_password: ENV["MONGO_INITDB_ROOT_PASSWORD"],
  jwt_encrypt_key: ENV["JWT_ENCRYPTION_KEY"],
  cookie_encrypt_key: ENV["COOKIE_ENCRYPTION_KEY"],
  client_id: ENV["AUTH0_CLIENTID"],
  client_secret: ENV["AUTH0_CLIENTSECRET"],

  authenticate_db_name: ENV["AUTH0_AUTH_DB"],
  authenticate_secret: ENV["AUTH0_CLIENTSECRET"],
  rabbit_username: ENV["RABBITMQ_DEFAULT_USER"],
  rabbit_password: ENV["RABBITMQ_DEFAULT_PASS"]
}

$command = {
  make: <<~CMD,
    docker run -d --name #{$config[:servicename]} --network infrastructure_network  --cap-add=IPC_LOCK --restart=on-failure:10 \
    -p #{$config[:host_port]}:#{$config[:docker_port]} \
    --volume #{$config[:vault_host_volume]}/:/vault/config \ #{$config[:image_name]} server
  CMD
  start: "docker start #{$config[:servicename]}",
  stop: "docker stop #{$config[:servicename]}",
  log: "docker logs -t #{$config[:servicename]}",
  prune: "docker rm -f #{$config[:servicename]}",
  stat: "docker stats #{$config[:servicename]}"
}

def vaultconfiggen
  erb_file_path = File.join(Dir.pwd, "script", "vault", "vault.config.erb")
  template = File.read(erb_file_path)
  vault_listen_address = $config[:vault_listen_address]
  vault_api_address = $config[:vault_api_address]
  secret_key = $config[:aws_secret_key]
  kms_key_id = $config[:aws_kms_key_id]
  access_key = $config[:aws_access_key]
  kms_region = $config[:aws_region]
  erb = ERB.new(template)
  result = erb.result(binding)
  file_path = File.join(Dir.pwd, "script", "vault-config", "vault-config.hcl")
  File.write(file_path, result)
end

def initVault
  ENV["VAULT_ADDR"] = $config[:vault_api_address]

  output = `vault operator init`
  recoveryKey = output.scan(/Recovery Key \d: (\S+)/).flatten

  recoveryKeyFile = File.join(Dir.pwd, "secret", "recovery_token.sample.txt")

  File.write(recoveryKeyFile, recoveryKey.join("\n"))

  rootToken = output.match(/Initial Root Token: (\S+)/)[1]

  rootTokenfile = File.join(Dir.pwd, "secret", "root_token.sample.txt")
  File.write(rootTokenfile, rootToken)

  ENV["VAULT_TOKEN"] = rootToken

  system("vault secrets enable -path=secret kv-v2")
  system("vault auth enable userpass")

  policy_name = "test_policy"
  policy_content = <<~EOF
    path "secret/web-server" {
      capabilities = ["read", "list"]
    }
    path "secret/data/web-server" {
      capabilities = ["read", "list"]
    }
    path "secret/web-server/cert" {
      capabilities = ["read" , "list"]
    }
    path "secret/data/web-server/cert" {
      capabilities = ["read" , "list"]
    }

    path "secret/web-server/auth" {
      capabilities = ["read" , "list"]
    }

    path "secret/data/web-server/auth" {
      capabilities = ["read" , "list"]
    }
  EOF

  command = <<~CMD
    vault policy write #{policy_name} - <<EOF
    #{policy_content.chomp}
    EOF
  CMD

  system(command)
  system("vault write auth/userpass/users/webserver password=password policies=#{policy_name}")
  system(
    "vault kv put secret/web-server MONGO_INITDB_ROOT_USERNAME=#{$config[:mongo_username]} MONGO_INITDB_ROOT_PASSWORD=#{$config[:mongo_password]} JWT_ENCRYPTION_KEY=#{$config[:jwt_encrypt_key]} COOKIE_ENCRYPTION_KEY=#{$config[:cookie_encrypt_key]} RABBIT_MQ_USERNAME=#{$config[:rabbit_username]} RABBIT_MQ_PASSWORD=#{$config[:rabbit_password]}"
  )
  system(
    "vault kv put secret/web-server/cert CERT=@secret/localhost.pem KEY=@secret/localhost-key.pem"
  )

  system(
    "vault kv put secret/web-server/auth CLIENTID=#{$config[:client_id]} CLIENTSECRET=#{$config[:authenticate_secret]} AUTH_DBNAME=#{$config[:authenticate_db_name]}"
  )
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
    vaultconfiggen
    system($command[command])
    sleep(2.5)
    initVault
  else
    system($command[command])
  end
end

main
