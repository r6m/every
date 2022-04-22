# every
Every translates english to crontab expression

| every expression     | crontab expression |
| ----------- | ----------- |
| every 2 minutes      | */2 * * * *       |
| every 2 minutes on Fri   | */2 * * * Fri        |
| every day at 3 pm | 0 15 * * * |
| every hour in Jun on Sun,Fri | 0 * * Jun Sun,Fri |


## TODO
- [ ] `every` cmd to run commands
- [ ] support Everyfile to generate crontab file (Caddyfile)
- [ ] read args to generate from input
