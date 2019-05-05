# clean db periodicaly

## notes for cross compile from Windows to Linux
```cmd
set GOARCH=amd64
set GOOS=linux
```

do `go env` to check env
for build `go build`

## Notes for Vault
file `.vault_pass` allows read password from env
put password into env
`export VAULT_PASSWORD=<password>`

view / edit vars
`ansible-vault view group_vars/ec2/vault.yaml`

#### tags
install / re-install only spiders service
`#### Notes for Vault
 file `.vault_pass` allows read password from env
 put password into env
`export VAULT_PASSWORD=<password>`
 
view / edit vars
`ansible-vault view group_vars/ec2/vault.yaml`
 
#### tags
install / re-install only spiders service
`ansible-playbook -t mongoCleaner deploy.yaml`
 

