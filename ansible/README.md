# Deploy Notes

TODO:
- add info about new host in `hosts.yaml`

## Notes for Vault
- file `.vault_pass` allows read password from env

- it need `x` permissions

to put password into env (space mate :-) )
` export VAULT_PASSWORD=<password>`

view / edit vars
`ansible-vault view group_vars/ec2/vault.yaml`

## tags
install / re-install only spiders service
`ansible-playbook -t mongoCleaner deploy.yaml`

## some manual job

### fresh OS

```
sudo adduser vic
```
```
~$ sudo grep vic /etc/passwd
vic:x:1001:1001:,,,:/home/vic:/bin/bash
```
```
~$ sudo grep vic /etc/group
adm:x:4:syslog,ubuntu,vic
dialout:x:20:ubuntu,vic
cdrom:x:24:ubuntu,vic
floppy:x:25:ubuntu,vic
sudo:x:27:ubuntu,vic
audio:x:29:ubuntu,vic
dip:x:30:ubuntu,vic
video:x:44:ubuntu,vic
plugdev:x:46:ubuntu,vic
lxd:x:108:ubuntu,vic
netdev:x:114:ubuntu,vic
vic:x:1001:vic
```
```
~$ sudo grep vic /etc/sudoers
vic ALL=(ALL) NOPASSWD:ALL
```

from ubuntu user 
```
sudo cp -R .ssh /home/vic/
sudo chown -R vic /home/vic/.ssh
sudo chgrp -R vic /home/vic/.ssh
```

### Migration to new instance

#### mongodump

- dump as is to dump dir
`mongodump --numParallelCollections=1`

- best option but, `fatal error: runtime: out of memory`
mongodump --numParallelCollections=1 --gzip --archive=`date +%Y-%m-%d`_mongodump

- tgz
```
tar cvzf `date +%Y-%m-%d`_mongodump.tgz dump
```

#### mongorestore

- `sudo systemctl status curs_auto`

- restore, with overwrite 
`mongorestore --drop --numParallelCollections=1 dump` Before restoring the collections from the dumped backup, drops the collections from the target database. 

