#!/bin/bash

swagger generate client -f teamcity-spec.yml -c client -a client -m models --config-file=template-config.yml --skip-validation