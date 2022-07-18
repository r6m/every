# every
Every translates english to crontab expression

| every expression     | crontab expression |
| ----------- | ----------- |
| every 2 minutes      | */2 * * * *       |
| every 2 minutes on Fri   | */2 * * * Fri        |
| every day at 3 pm | 0 15 * * * |
| every hour in Jun on Sun,Fri | 0 * * Jun Sun,Fri |

## Everyfile config
**every** uses hashicorp hcl config structure. 
```hcl
every "2 minutes" {
  run = "command >/dev/null 2>&1"
  user = "ubuntu"
}

every "day at 3 pm" {
  run = "command >/dev/null 2>&1"
  user = "ubuntu"
}

every "hour in Jun on Sun,Fri" {
  run = "command >/dev/null 2>&1"
  user = "ubuntu"
}
```

## TODO
- [x] `every` cmd to run commands
- [x] support Everyfile to generate crontab file (hcl)
- [ ] read args to generate from input
- [ ] run jobs as specific user
- [x] update crontab from config
