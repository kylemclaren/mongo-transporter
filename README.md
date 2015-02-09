![Mongo Transporter](mongo_transporter.png)

[![wercker status](https://app.wercker.com/status/3eda307e1ccd93047fb764846c90bc9b/m/master "wercker status")](https://app.wercker.com/project/bykey/3eda307e1ccd93047fb764846c90bc9b)

A simple Go app that uses the [Compose Transporter](https://github.com/compose/transporter) to transfer data between two MongoDB deployments and keep them in sync.

## What's it good for?

- Keeping dev, staging and production DB's in sync.
- Zero-downtime migrations from one deployment to another.

<!--

## What it does

- connect to both the source and the destination and finds the oplog timestamp
- copies unique indexes from source to destination (changing their namespace)
- copies users
- copies all the collections in parallel
- copies non-unique indexes
- tails the oplog from the initial timestamp, and applies the operations in a batch (ignoring a list of blacklisted - - commands, dropDatabase, etc). There is no conflict resolution with Transporter. When writing to the source and the destination, the last write always wins.

-->

## Deploy!

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/kylemclaren/mongo-transporter)

Click the deploy button to launch a new app instance, add your config/environment variables in the Heroku dashboard and click "Deploy for Free". This will create a new Heroku app. Worker dynos do not scale automatically when deploying, you will need to scale manually to one worker dyno via the dashboard or the command line: `heroku ps:scale worker=1 -a 'YOUR APP NAME'`

For now, Mongo Transporter will only sync a single DB on a deployment so if a deployment has multiple DB's, you will have to run multiple app instances.

## Config vars

- `SOURCE_MONGO_URL` - This is the full connection URI of the MongoDB deployment that you want to sync **from**. eg. `mongodb://username:strongpassword@candidate.33.mongolayer.com:30000,candidate.34.mongolayer.com:30000/local?authSource=prod_db` You will need to create a user that can read from the `local.oplog.rs` namespace. You can use both members if the deployment is a replica set.
- `SOURCE_DB` - The DB name to sync from. eg. `prod_db`
- `SINK_MONGO_URL` - This is the full connection URI of the MongoDB deployment that you want to sync **to**. eg. `mongodb://username:strongpassword@candidate.43.mongolayer.com:30000,candidate.44.mongolayer.com:30000/staging_db` The user does not need to authenticate to the `local` DB but needs read write access to `SINK_DB`. You can use both members if the deployment is a replica set.
- `SINK_DB` - The DB name to sync from. eg. `staging_db`
- `TAIL` - Specify true to run a continuous sync, tailing the oplog. False for a one-time sync.
- `DEBUG` - Specify true for verbose logging to stdout.

Run `$ heroku logs -ta MY_APP_NAME` from the command line to check out the logs.

<!-- Note that the users for both the source and destination deployments must use a user with [oplog access](https://docs.compose.io/common-questions/getting-oplog-access.html). -->

## What is does

- Connect to both the source and the destination and finds the oplog timestamp
- Copies all the collections in parallel
- Tails the oplog from the initial timestamp, and applies the operations in a batch (ignoring a list of blacklisted commands - dropDatabase, etc). There is no conflict resolution with Transporter. When writing to the source and the destination, the last write always wins.

## What it does not do

- Copy DB users
- Copy indexes

You will have to recreate any (non-unique) indexes and create new users on the destination DB, a small trade-off for the ease of use.

## Thanks

The engineers at [Compose](https://compose.io) for making an awesome tool, Transporter.


This app uses the [Go Buildpack for Heroku](https://github.com/kr/heroku-buildpack-go) by @kr

## To Do

- [ ] 2.4/2.6 caveats (if any)
- [ ] Comma separated list of collections to ignore
- [ ] Improve log output
- [ ] Slack notifications

## License

Copyright (c) 2014, Compose, Inc.

All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of [project] nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
