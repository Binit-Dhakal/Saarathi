#!/bin/bash

buf generate

buf generate --template dto/buf.gen.yaml --path shared/common --path dto/driverspb
