#!/bin/bash
minishift profile set minishift
minishift profile delete $CLUSTER_UUID -f
