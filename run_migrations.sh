#!/bin/sh
migrate -path ./migrations -database "mysql://root:example@tcp(localhost:3307)/todoapp" up