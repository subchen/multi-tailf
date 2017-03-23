# multi-tailf
watch multiple logs on local or remote servers.

# Dependencies

- sshpass: http://sourceforge.net/projects/sshpass/

# Usage

```
Usage: mtailf [OPTIONS] [ /path/file | user:pass@host:file ] ...
   or: mtailf [ --version | --help ]

mtailf is equivalent of tail -f on multiple local or remote files at once

Options:
  --ssh-pass=value          default password for ssh
  --ssh-key=~/.ssh/id_rsa   default key file for ssh
  --ssh-file=value          default file for ssh tail
  --version                 show version information
  --help                    show this help

Examples:
* Local files
    mtailf /var/log/messages-1 /var/log/messages-2
* Multiple files on servers
    mtailf root@10.0.0.1 root@10.0.0.2 --ssh-pass=password --ssh-file=/var/log/messages
* Use SSH private key
    mtailf root@10.0.0.1:/var/log/messages --ssh-key=/tmp/ssh.key
* Use SSH passwords
    mtailf root:p1@10.0.0.1:/var/log/messages root:p2@10.0.0.2:/var/log/messages
* Use SSH port
    mtailf root@10.0.0.1:8022:/var/log/messages
```
