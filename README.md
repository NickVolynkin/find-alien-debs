# Find alien debs

This tool finds packages which are absent in locally cached package lists. So,
you may use it for inspect your system for manually-installed packages or some kind
of distro version mixup.

E.g. you have enabled `jessie-backports` and `jessie`, but have some fresh and broken
libs from `stretch`. Just disable `stretch` and...

```
$ sudo apt-get update
$ find-alien-debs | grep absent
flake8:all 2.5.4-3 absent
libfdisk1:amd64 2.28-5 absent
tlp:all 0.8-1 absent
xnview:amd64 0.79 absent
```

Requirements: Debian-based distro

