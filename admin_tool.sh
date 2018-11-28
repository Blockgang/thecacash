#!/bin/bash
ENV_FILE=.env
source ${ENV_FILE}

stackCommonServicesFile=docker-compose.stack_common-services.yml
stackBuildFile=docker-compose.stack_build.yml
stackRunFile=docker-compose.stack_run.yml

options=(
    "create_docker_secrets" \
    "build" \
    "build_hard" \
    "run" \
    "run_hard" \
    "removeStacks" \
    "removeContainers" \
    "show_logs_build" \
    "show_logs_run" \
    "Quit"
)

function createComposeFile {
    source ${ENV_FILE}
    buildSecretConfig=docker-compose.build_template.yml
    runSecretConfig=docker-compose.run_template.yml
    docker-compose -f docker-compose.build_template.yml -f ${buildSecretConfig} config > ${stackBuildFile}
    docker-compose -f docker-compose.run_template.yml -f ${runSecretConfig} config > ${stackRunFile}
}

function evalParams () {
    containers=""
    while (( "$#" )); do
        if [ "$1" == "hard" ]; then
            buildHardOptions="--force-rm --no-cache --pull"
        elif [ "$1" == "build" ]; then
            composeFile=${stackBuildFile}
        elif [ "$1" == "run" ]; then
            composeFile=${stackRunFile}
        else
            containers+=$1" "
        fi
        shift
    done
}

function buildContainer () {
    evalParams $@
    docker-compose -f ${composeFile} build ${buildHardOptions} ${containers}
}

function startContainer () {
    evalParams $@
    source ${ENV_FILE}
    docker-compose -f ${composeFile} up -d $containers
}

function removeStacks () {
    docker stack rm ${DOCKER_STACK_NAME_BUILD}
    docker stack rm ${DOCKER_STACK_NAME_RUN}
}

function removeContainers () {
    createComposeFile compose
    docker-compose -f ${stackBuildFile} -f ${stackRunFile} rm -sf
}

function showLogs () {
    evalParams $@
    source ${ENV_FILE}
    docker-compose -f ${composeFile} logs -f $containers
}

function doit {
    case "$1" in
        "create_docker_secrets")
            echo "create_docker_secrets"
            ./create_docker_secrets.sh
            ;;
        "build")
            echo "build theca.cash ${*:2}"
            createComposeFile
            buildContainer build ${*:2}
            startContainer build ${*:2}
            ;;
        "build_hard")
            echo "build theca.cash hard ${*:2}"
            createComposeFile
            buildContainer build hard ${*:2}
            docker stack rm ${DOCKER_STACK_NAME_BUILD} 2> /dev/null
            # sleeping a while to delete docker stack
            sleep 20
	    startContainer build hard ${*:2}
            ;;
        "run")
            echo "run theca.cash ${*:2}"
            createComposeFile
            buildContainer run ${*:2}
            startContainer run ${*:2}
            ;;
        "run_hard")
            echo "run theca.cash hard ${*:2}"
            createComposeFile
            buildContainer run hard ${*:2}
	    docker stack rm ${DOCKER_STACK_NAME_RUN} 2> /dev/null
            # sleeping a while to delete docker stack
            sleep 20
            startContainer run hard ${*:2}
            ;;
        "removeContainers")
            echo "removeContainers"
            removeContainers
            ;;
        "show_logs_build")
            echo "show logs build dockers ${*:2}"
            showLogs build ${*:2}
            ;;
        "show_logs_run")
            echo "show logs run dockers ${*:2}"
            showLogs run ${*:2}
            ;;
        "Quit")
            echo "exit"
            exit 0
            ;;
        *)
            echo "unknown option"
            ;;
    esac
}

function print_usage {
    echo -e "Usage: admin_tool.sh [option] [OPTIONAL_DOCKER_TO_BUILD (default: all) docker-service]\n" \
            "option:  -b        to build containers and build theca.cash in docker\n" \
            "         -bh       to build containers and build theca.cash in docker (build --force-rm --no-cache --pull)\n" \
            "         -r        to build containers and run theca.cash in docker\n" \
            "         -rh       to build containers and run theca.cash in docker (build --force-rm --no-cache --pull) \n" \
            "         -lb       to docker service logs from build dockers\n" \
            "         -lr       to docker service logs from run dockers\n"
}

if (( $# == 0 ))
then
    if (( $OPTIND == 1 ))
    then
        PS3='Please enter your choice: '
        select opt in "${options[@]}"
        do
            doit "${opt}"
        done
    fi
fi

case "$1" in
    "-b")
        doit build ${*:2}
        ;;
    "-bh")
        doit build_hard ${*:2}
        ;;
    "-r")
        doit run ${*:2}
        ;;
    "-rh")
        doit run_hard ${*:2}
        ;;
    "-lb")
        doit show_logs_build ${*:2}
        ;;
    "-lr")
        doit show_logs_run ${*:2} 
        ;;
    "-h")
        print_usage
        exit 0
        ;;
    *)
        echo -e "argument error \n"
        print_usage
        exit 1
        ;;
esac

