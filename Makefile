# Makefile CHEATSHEET: https://devhints.io/makefile
##############################################################################
include Makefile.settings
##############################################################################
# Meta

menu :
	$(INFO) 'Test source:'
	@echo '	test  : go test ./...'
	$(INFO) 'Run CLI (per makeargs):'
	@echo '	gorun  : go run ./app/cli ...'
	$(INFO) 'Push to repo:'
	@echo '	push  : git push -u origin master'

env :
	@env |grep APP_

##############################################################################
# Source 

test :
	go test ./...

# git remote add origin git@github.com:$_USERNAME/$_REPONAME.git  # ssh mode
push :
	git push -u origin master

markup :
	bash make.md2html.sh

perms :
	bash make.perms.sh

##############################################################################
# App [CLI] : Set `makeargs` to the command plus its options (default: env)

# USAGE: 
# 	export makeargs='..'
# 	make gorun

gorun :
	bash make.go.run.app.sh cli $(shell echo "$${makeargs:-env}")

