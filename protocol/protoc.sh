#!/bin/bash
protoc \
-I=./proto/ \
--go_out=./proto/gate/ \
--go_opt=paths=source_relative \
./proto/gate.proto


protoc \
-I=./proto/ \
--go_out=./proto/login/ \
--go_opt=paths=source_relative \
./proto/login.proto


protoc \
-I=./proto/ \
--go_out=./proto/common/id \
--go_opt=paths=source_relative \
./proto/errorId.proto


protoc \
-I=./proto/ \
--go_out=./proto/common/id \
--go_opt=paths=source_relative \
./proto/messageId.proto


protoc \
-I=./proto/ \
--go_out=./proto/common/content \
--go_opt=paths=source_relative \
./proto/content.proto


protoc \
-I=./proto/ \
--go_out=./proto/server/cluster \
--go_opt=paths=source_relative \
./proto/cluster.proto