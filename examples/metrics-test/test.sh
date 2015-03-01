#!/bin/bash

curlit() {
    header="$1"
    curl -H "Content-Type: $header" http://localhost:8111
}

curlit 'application/json'
curlit 'application/xml'
curlit 'application/html'
curlit 'unknown'
