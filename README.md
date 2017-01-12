# Simple Reminder Telegram Bot 

## Simple golang and mongodb driven reminder bot

### Requirements

Reminder Bot version 0.1 requires Go >= 1.7 and MongoDB >=3.2.x

##### Installation

```sh
$ go get github.com/alexivanenko/my_simple_reminder_bot/...
$ cd my_simple_reminder_bot
$ cp config.sample config.ini
$ make
```

Create 'simple_reminder' DB in the MongoDB installation, 'events' collection and also add the new DB user.
Mongo authentication credentials in the config.ini file.

Next step is installation of the notifications manager.

```sh
$ go get github.com/alexivanenko/my_simple_reminder_manager/...
$ cd my_simple_reminder_manager
$ cp config.sample config.ini
$ make
```
