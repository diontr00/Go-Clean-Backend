listener "tcp" { 
  address =  "0.0.0.0:8400"
  tls_disable = true 
}

backend "file" { 
  path = "vault/file"
}

seal "awskms" { 
  region =  "ap-southeast-2"
  kms_key_id = "<your_kms_key_id>"
  secret_key = "<your_aws_secret_key>"
  access_key = "<your_aws_access_key>"
}


api_addr = "http://0.0.0.0:8400"
ui = true 
