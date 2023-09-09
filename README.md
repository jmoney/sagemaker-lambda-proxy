# Sagemaker Lambda Proxy

A repo for the terraform required to standup a sagemaker proxy using lambdas in AWS

## Prerequisites

```bash
brew bundle --file Brewfile
tfenv use
```

## Deploy

```bash
make build
terraform init
AWS_PROFILE=<profile> terraform plan -out=plan.out -var endpoint_name=<endpoint-name>
AWS_PROFILE=<profile> terraform apply plan.out
```

## Destroy

```bash
AWS_PROFILE=<profile> terraform destroy -var endpoint_name=<endpoint-name>
```

## Testing

```bash
curl --silent --http1.1 --request POST --header "Content-Type: application/json" --header "X-NONCE: $(terraform output -json | jq -r .nonce.value)" "$(terraform output -json | jq -r .url.value)" -d @data.json | jq -r '.[].generation.content'
```

There is a `data.json` file at the root of this repo with a question you can use.

## Docker

There is a dockerfile built to be used with lambda for the container image package type.  It has not been tested for local runs but in theory it should work.  Lambda does not allow images pulled from ghcr and only AWS ECR registries so the docker image would need to be mirrored from ghcr to ecr in some way.
