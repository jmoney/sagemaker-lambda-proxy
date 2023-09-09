# llama2-terraform

A repo for the terraform required to standup llama2 proxy in AWS

## Prerequisites

```bash
brew bundle --file Brewfile
tfenv use
```

## Deploy

```bash
make build
terraform init
AWS_PROFILE=<profile> terraform plan -out=plan.out -var endpoint-name=<endpoint-name>
AWS_PROFILE=<profile> terraform apply plan.out
```

## Destroy

```bash
AWS_PROFILE=<profile> terraform destroy
```

## Testing

```bash
curl --silent --method POST --header "NONCE: $(terraform output -json | jq -r .nonce.value)" --data @data.json "$(terraform output -json | jq -r .url.value)"
```
