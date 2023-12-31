# JetBrains AI Test Drive

This repo contains code that was generated by the JetBrains AI in Goland in December 2023, while I was on PTO for the Cisco holiday shutdown.

THe goal of this project was to see how much code I could prompt the JetBrains AI (the JBAI) into generating for me, and how close I could get to something semi-releasable. This is not a very complex codebase, as I wanted to start off slow and maybe work my way up from here.

The full story behind it can be seen in [this blog series](https://colinj.hashnode.dev/series/chat-driven-development).

## Goals:

- build a layered-architecture microservice using prompts
- get as close to 100% unit test coverage as possible with the majority of test code written by the JBAI
- as little manual code as possible
- it has to work - that is, runnable code that can receive and respond to REST requests

### Tradeoffs:

- data storage is mocked with an in-memory map in order to limit the scope of the project
- code is kept as-is if possible, to preserve AI weirdness, though corrections were made to make sure the code runs
- project structure does not entirely follow best practices, for the purposes of keeping things simple