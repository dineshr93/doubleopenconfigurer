
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
PROJ_DIR := ${ROOT_DIR}proj
BKP_DIR := bkp
EXE_NAME := parseconf
BIN := bin

ifeq ($(OS),Windows_NT)
	SHELL := powershell.exe
	.SHELLFLAGS := -NoProfile -Command
	RM_F_CMD = Remove-Item -erroraction silentlycontinue -Force
    RM_RF_CMD = ${RM_F_CMD} -Recurse
	exe =${BIN}/${EXE_NAME}.exe
	HELP_CMD = Select-String "^[a-zA-Z_-]+:.*?\#\# .*$$" "./Makefile" | Foreach-Object { $$_data = $$_.matches -split ":.*?\#\# "; $$obj = New-Object PSCustomObject; Add-Member -InputObject $$obj -NotePropertyName ('Command') -NotePropertyValue $$_data[0]; Add-Member -InputObject $$obj -NotePropertyName ('Description') -NotePropertyValue $$_data[1]; $$obj } | Format-Table -HideTableHeaders @{Expression={ $$e = [char]27; "$$e[36m$$($$_.Command)$${e}[0m" }}, Description
else
	SHELL := bash
	RM_F_CMD = rm -f
	RM_RF_CMD = ${RM_F_CMD} -r
	exe =${BIN}/${EXE_NAME}
	HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
endif


.DEFAULT_GOAL := help
.PHONY: clean build test all git help

all: clean build test ## performs clean build and test


build: ## Generate the windows and linux builds for sep
	echo "Compiling for every OS and Platform"
	set GOOS=windows
	set GOARCH=arm64
	set CGO_ENABLED=0
	CGO_ENABLED=0 go build -o ${BIN}/${EXE_NAME}.exe doubleopenconf.go utils.go
	set GOOS=linux
	set GOARCH=amd64
	CGO_ENABLED=0 go build -o ${BIN}/${EXE_NAME} doubleopenconf.go utils.go


test: ## downloads the meta-doubleopen git repo and configures 2 files in test folder workdir/conf/ bblayers.conf and local.conf
	echo "===========Testing==============="
	${exe} ${PROJ_DIR} 18 14

del: ## Delete the contents of test (named: proj) folder
	${RM_RF_CMD} ${PROJ_DIR}/*
	${RM_RF_CMD} ${BIN}/*

cp:  ## Copy Sample conf files to test (named: proj) folder in workdir/conf/
	mkdir ${PROJ_DIR}/workdir
	mkdir ${PROJ_DIR}/workdir/conf
	cp ${BKP_DIR}/* ${PROJ_DIR}/workdir/conf


clean: del cp ## Copy back the files to original state in test folder


git: ## commits and push the changes if commit msg m is given without spaces ex m=added_files
	git status
	git add .
	git status
	git commit -m ${m}
	git push

help: ## Show this help
	@${HELP_CMD}
