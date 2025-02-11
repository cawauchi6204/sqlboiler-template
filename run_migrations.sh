#!/bin/sh
migrate -path ./migrations -database "mysql://root:example@tcp(db:3306)/todoapp" up