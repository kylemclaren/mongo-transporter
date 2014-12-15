<!-- ![Mongo Transporter](mongo_transporter.png) -->

A dead simple Go app that uses the [Compose Transporter](https://github.com/compose/transporter) to transfer data between two MongoDB deployments and keep them in sync.

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

Launch a new app and then add your environment variables in the Heroku dashboard. This will create a new Heroku app. You will need to scale manually to one worker dyno via the dashboard or the command line `heroku ps:scale worker=1`

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/kylemclaren/mongo-transporter)

## Config vars

- `SOURCE_MONGO_URL` - This is the full connection URI of the MongoDB deployment that you want to sync **from**. eg. `mongodb://username:strongpassword@candidate.44.mongolayer.com:10000/local?authSource=prod_db` You will need to create a user that can read from the `oplog.rs` collection.
- `SOURCE_NS` - The DB name and collection to sync from. eg. `prod_db.myCollection`
- `DESTINATION_MONGO_URL` - This is the full connection URI of the MongoDB deployment that you want to sync **to**. eg. `mongodb://username:strongpassword@candidate.44.mongolayer.com:30000/staging_db?authSource=local` This must be the Primary member of the replica set and can use any user with read\write access. THe user doesn not need to authenticate to the `local` DB.
- `DEST_NS` - The DB name and collection to sync from. eg. `staging_db.myCollection`

Note that the users for both the source and destination deployments must use a user with [oplog access](https://docs.compose.io/common-questions/getting-oplog-access.html).

## Thanks

This app uses the Go Buildpack for Heroku by @kr: https://github.com/kr/heroku-buildpack-go

## To Do

1. 2.4/2.6 caveats

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
