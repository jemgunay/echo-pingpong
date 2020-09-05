# Ping Pong Counter Alexa Skill Server

An Amazon Alexa skill for keeping track of ping pong game scores.

## Deployment

Deploy to Google cloud (will require account setup prior):
```bash
make deploy
```

Stream logs:
```bash
make attach_log
```

## Example Chatter Sequence

* [C] Play a new ping pong game
* [C] Add player {name 1}
* [S] {name 1} added
* [C] Add player {name 2}
* [S] {name 2} added
* [C] Start game
* [S] Player {random, name 1} to serve
* [C] {name 1} scored
* [S] 1:0 to {name 1}. {name 1} to serve
* [C] Point to {name 2}
* [S] 1 all. {name 2} to serve
...
* [C] {name 2} scored
* [S] {name 2} wins. Play again?

## TODO

* Match point dialog
* "What's the score?" action
* Ability to set up tournaments with > 2 people