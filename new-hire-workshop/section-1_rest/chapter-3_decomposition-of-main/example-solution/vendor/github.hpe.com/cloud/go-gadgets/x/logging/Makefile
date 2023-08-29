# (C) Copyright 2021-2022 Hewlett Packard Enterprise Development LP

# This import must be present for the root level Makefile to function
include ../../module.mk

# NOTE: This module uses 'golden' files to assert the correct output of loggers
# in some of its tests.
# Golden files are sample output files that are used to compare against in unit tests.
# If an implementation has been deliberately changed, then the expected output in the golden
# file will need to be updated. Manual copy/paste of output is prone to error (e.g. spacing),
# and therefore a better solution is to allow the test to write out the expected output
# to the file.

## update-golden: Update the golden files with the current output of the unit tests.
.PHONY: update-golden
update-golden: download
	$(GOTEST) ./... \
		-update
