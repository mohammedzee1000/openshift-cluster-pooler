#!/bin/bash

minishift profile set $CLUSTER_UUID &> /dev/null
minishift_ip=$(minishift ip)
echo "https://$minishift_ip:8443"