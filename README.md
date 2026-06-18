# honeypot

An old trick. Real SSH on port 22 just attracts noise — bots, scanners, the
usual midnight traffic. So put something on port 22 that isn't a shell.

This is a tiny SSH server that accepts any connection, writes a single line, and
hangs up. No auth prompt, no shell, nothing to brute-force. Whoever knocks gets
a message and a closed door.

## Run

```
go run main.go
```

Port 22 needs root. The line it prints comes from `HONEYPOT_EMAIL` in a `.env`
file — copy `.env.example` to `.env` and set your address before you start.

## How it holds up

Built to sit untouched on a box for a long time:

- one token bucket per source IP, so a single host can't flood it
- pairs with a `fail2ban` jail (`honeypot-ssh`) for repeat offenders
- connections logged to `/var/log/honeypot.log`

Nothing exciting happens, which is the point. You move the real SSH to another
port, leave this one running, and forget about it for a few years.
