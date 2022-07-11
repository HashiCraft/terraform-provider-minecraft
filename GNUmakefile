default: testacc

# Run acceptance tests
.PHONY: testacc example
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

name = minecraft
organization = hashicraft
version = 0.1.0
log_level = info

build:
	go build -o bin/terraform-provider-$(name)_v$(version)

install: build clean
	mkdir -p ~/.terraform.d/plugins/local/$(organization)/$(name)/$(version)/linux_amd64
	mv bin/terraform-provider-$(name)_v$(version) ~/.terraform.d/plugins/local/$(organization)/$(name)/$(version)/linux_amd64/

clean:
	rm -rf example/.terraform*
	rm -rf example/terraform.tfstate*

init:
	TF_LOG=$(log_level) terraform -chdir=examples init

plan:
	TF_LOG=$(log_level) terraform -chdir=examples plan

apply:
	TF_LOG=$(log_level) terraform -chdir=examples apply -auto-approve

destroy:
	TF_LOG=$(log_level) terraform -chdir=examples destroy -auto-approve