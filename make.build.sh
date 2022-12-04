#!/usr/bin/env bash
#------------------------------------------------------------------------------
#  Makefile : make build* SVC VER : docker build ... per environment 
# 
#  ARGs: SVC VER
# -----------------------------------------------------------------------------

export image="app.${1}.dockerfile"

docker build \
    -f ${APP_INFRA}/docker/build/$image \
    -t ${HUB}/${PRJ}.${1}-${ARCH}:${2} \
    --build-arg PKG_NAME=$1 \
    --build-arg ARCH=${ARCH} \
    --build-arg HUB=${HUB} \
    --build-arg PRJ=${PRJ} \
    --build-arg MODULE=${MODULE} \
    --build-arg AUTHORS="${AUTHORS}" \
    --build-arg VENDOR="${VENDOR}" \
    --build-arg SVN=${SVN} \
    --build-arg VER=${2} \
    --build-arg BUILT="${BUILT}" \
    . && docker push ${HUB}/${PRJ}.${1}-${ARCH}:${2}

exit

# Docker build options  https://docs.docker.com/engine/reference/commandline/build/#options
    --cpuset-cpus=2 \
    --cpuset-mems=512m \
    --memory=2g \
    --shm-size=256m  # /tmpfs # NO EFFECT
    --ulimit
    --iidfile

# @ adm : Nominal
/app $ df -h
Filesystem                Size      Used Available Use% Mounted on
overlay                  23.5G      3.5G     18.8G  16% /
tmpfs                    64.0M         0     64.0M   0% /dev
tmpfs                   994.4M         0    994.4M   0% /sys/fs/cgroup
shm                      64.0M         0     64.0M   0% /dev/shm
grpcfuse                 90.0G     36.9G     53.1G  41% /app/assets
/dev/sda1                23.5G      3.5G     18.8G  16% /etc/resolv.conf
/dev/sda1                23.5G      3.5G     18.8G  16% /etc/hostname
/dev/sda1                23.5G      3.5G     18.8G  16% /etc/hosts
tmpfs                   994.4M      4.0K    994.4M   0% /run/secrets/pg_pw_app_user
tmpfs                   994.4M         0    994.4M   0% /proc/acpi
tmpfs                    64.0M         0     64.0M   0% /proc/kcore
tmpfs                    64.0M         0     64.0M   0% /proc/keys
tmpfs                    64.0M         0     64.0M   0% /proc/timer_list
tmpfs                    64.0M         0     64.0M   0% /proc/sched_debug
tmpfs                   994.4M         0    994.4M   0% /sys/firmware

