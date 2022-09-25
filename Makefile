# Makefile CHEATSHEET: https://devhints.io/makefile
##############################################################################
include Makefile.settings
##############################################################################
# Meta

menu :
	$(INFO) 'Test source:'
	@echo '	test  : go test ./…'
	$(INFO) 'Run CLI (per makeargs):'
	@echo '	gorun  : go run ./app/cli …'
	$(INFO) 'Push to repo:'
	@echo '	push  : git push -u origin master'
	@echo '	tag   : git tag v${VER_APP}  (VER_APP)'
	@echo '	untag : git … : remove v${VER_APP}  (VER_APP)'
	$(INFO) 'Demo :'
	@echo '	token  : go run ./app/cli token'
	@echo '	upsert : go run ./app/cli upsert $$json_body $$mid $$tkn $$APP_CHANNEL_SLUG'

env :
	@env |grep APP_

##############################################################################
# Source 

test :
	go test ./...

# git remote add origin git@github.com:$_USERNAME/$_REPONAME.git  # ssh mode
push :
	gc
	git push -u origin master
tag :
ifeq (v${VER_APP}, $(shell git tag |grep v${VER_APP}))
	@echo 'repo ALREADY tagged @ "v${VER_APP}" : VER_APP'
else 
	git tag v${VER_APP}
	git push origin v${VER_APP}
	git tag
endif
untag :
	git tag -d v${VER_APP}
	git push origin --delete v${VER_APP}
markup :
	bash make.md2html.sh
tarball :
	bash make.tarball.sh
perms :
	bash make.perms.sh
tidy :
	go mod tidy
	go mod vendor

##############################################################################
# App [CLI] : Set `makeargs` to the command plus its options (default: env)

# USAGE: 
# ☩ export makeargs='..'
# ☩ make gorun

gorun :
	bash make.go.run.app.sh cli $(shell echo "$${makeargs:-env}")

# ☩ make APP_SERVICE_BASE_URL=https://uqrate.org goruntoken
goruntoken :
	bash make.go.run.app.sh token
gorunupsert :
	bash make.go.run.app.sh upsert