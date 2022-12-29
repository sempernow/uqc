# Makefile CHEATSHEET: https://devhints.io/makefile
##############################################################################
include Makefile.settings
##############################################################################
# Meta

menu :
	$(INFO) 'Manage source code :'
	@echo '	tidy      : go mod tidy;go mod vendor'
	@echo '	pkglist   : go list …'
	@echo '	test  : go test ./…'
	@echo '	push  : git push -u origin master'
	@echo '	tag   : git tag v${VER_APP}  (VER_APP)'
	@echo '	untag : git … : remove v${VER_APP}  (VER_APP)'

	$(INFO) 'Docker Build / Stack :'
	@echo '	pull   : docker pull …'
	@echo '	build  : docker build …; docker push … '
	@echo '	prune  : docker image prune -f'
	@echo '	up     : docker stack deploy … '
	@echo '	down   : docker stack rm … '
	@echo '	new    : docker service logs … |grep "HTTP 201"'
	@echo '	cli    : docker exec -it … '

	$(INFO) 'Run CLI @ Demos :'
	@echo '	token  : go run ./app/cli token'
	@echo '	uptkn  : go run ./app/cli uptkn $$json $$mid $$tkn $$APP_CHANNEL_SLUG'
	@echo '	upkey  : go run ./app/cli upkey $$json $$mid $$key'

	$(INFO) 'Run CLI : Configure per `export makeargs="cli …"` :'
	@echo '	gorun   : go run ./app/cli $$makeargs'

	$(INFO) 'Run CLI : Preconfigured to process WordPress sites list : `go gun …` :'
	@echo '	siteslist   : Make a new sites list from CSV sources list.'
	@echo '	updateusers : Update all users of sites list.'
	@echo '	upsertchns  : Upsert all channels of sites list.'
	@echo '	upsertposts : Upsert all posts of all sites on sites list.'

	$(INFO) 'Purge cache :'
	@echo '	purgecachetkns  : rm ${APP_CACHE}/tkn.*'
	@echo '	purgecacheposts : rm ${APP_CACHE}/*_posts.json'

	$(INFO) 'Current operational settings : The "*" indicates it affects others :'
	@echo '	Hypervisor *    HYPERVISOR          : ${HYPERVISOR} (os|hyperv|aws)'
	@echo '	ServiceHost     APP_SERVICE_HOST    : ${APP_SERVICE_HOST}'
	@echo '	ClientUser      APP_CLIENT_USER     : ${APP_CLIENT_USER}'
	@echo '	SitesListCSV    APP_SITES_LIST_CSV  : ${APP_SITES_LIST_CSV}'
	@echo '	SitesListJSON   APP_SITES_LIST_JSON : ${APP_SITES_LIST_JSON}'

env :
	@env |grep APP_

##############################################################################
# Source 

test :
	go test ./...

pkglist : 
	$(INFO) '/app'
	@push ./app; go list -f '{{ .Name | printf "%14s" }}  {{ .Doc }}' ./...;pop
	$(INFO) '/client'
	@push ./client; go list -f '{{ .Name | printf "%14s" }}  {{ .Doc }}' ./...;pop
	$(INFO) '/kit'
	@push ./kit; go list -f '{{ .Name | printf "%14s" }}  {{ .Doc }}' ./...;pop

pull : 
	docker pull ${HUB}/${PRJ}.cli-${ARCH}:${VER_APP}
fixdocker credstore:
	sed -i 's/credsStore/credStore/g' ~/.docker/config.json

build : fixdocker
	bash ./make.build.sh cli ${VER_APP}
up :
	docker stack deploy -c ${APP_INFRA}/docker/services/stack-adm.yml ${PRJ}
down :
	docker stack rm ${PRJ}
new :
	docker service logs uqc_cli 2>&1 |grep "HTTP 201"
cli : 
	bash ./make.exec.sh cli

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
# ☩ export makeargs='cli trace https://uqrate.org/liveness'
# ☩ make gorun |jq .

gorun :
	@bash make.go.run.app.sh $${makeargs:-cli}

# ☩ make APP_SERVICE_BASE_URL=https://uqrate.org goruntoken
token tkn :
	@bash make.go.run.app.sh token
key :
	@bash make.go.run.app.sh key

uptkn :
	@bash make.go.run.app.sh uptkn

upkey :
	@bash make.go.run.app.sh upkey

site :
	@go run ./app/cli site

siteslist :
	@bash make.go.run.app.sh siteslist

updateusers :
	@bash make.go.run.app.sh updateusers
upsertchns :
	@bash make.go.run.app.sh upsertchns
upsertposts :
	@bash make.go.run.app.sh upsertposts

purgecachetkns :
	@bash make.go.run.app.sh purgecachetkns
purgecacheposts :
	@bash make.go.run.app.sh purgecacheposts

prune : 
	docker image prune -f
