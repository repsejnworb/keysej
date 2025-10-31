# keysej sshconf
## create a rule (guards that key exists and is keysej)
keysej sshconf new github github.com

## create a CIDR rule
keysej sshconf new work 10.16.0.0/16 --user jesper --forward

## list
keysej sshconf list work
keysej sshconf list --host github.com

## validate all keysej config fragments
keysej sshconf validate

## tidy files (sort, whitespace)
keysej sshconf tidy

## delete (with backup)
keysej sshconf delete work
keysej sshconf delete work 10.16.0.0/16