#!/usr/bin/python
# -*- coding: utf-8 -*-

import commands

ENV_FILE = ".env"
# source ${ENV_FILE}

stackCommonServicesFile = "docker-compose.stack_common-services.yml"
stackBuildFile = "docker-compose.stack_build.yml"
stackRunFile = "docker-compose.stack_run.yml"


def print_usage():
    print("Usage: admin_tool.sh [option] [OPTIONAL_DOCKER_TO_BUILD (default: all) docker-service]\n" \
            "option:  -b        to build containers and build theca.cash in docker\n" \
            "         -bh       to build containers and build theca.cash in docker (build --force-rm --no-cache --pull)\n" \
            "         -r        to build containers and run theca.cash in docker\n" \
            "         -rh       to build containers and run theca.cash in docker (build --force-rm --no-cache --pull) \n" \
            "         -lb       to docker service logs from build dockers\n" \
            "         -lr       to docker service logs from run dockers\n")

def createComposeFile():
    # source ${ENV_FILE}
    buildSecretConfig = docker-compose.build_template.yml
    runSecretConfig = docker-compose.run_template.yml
    commands.getoutput("docker-compose -f docker-compose.build_template.yml -f ${buildSecretConfig} config > ${stackBuildFile}")
    commands.getoutput("docker-compose -f docker-compose.run_template.yml -f ${runSecretConfig} config > ${stackRunFile}")

def removeStacks():
    commands.getoutput("docker stack rm ${DOCKER_STACK_NAME_BUILD}")
    commands.getoutput("docker stack rm ${DOCKER_STACK_NAME_RUN}")

def removeContainers():
    commands.getoutput("createComposeFile compose")
    commands.getoutput("docker-compose -f ${stackBuildFile} -f ${stackRunFile} rm -sf")

def createDockerNetwork():
    # source ${ENV_FILE}
    commands.getoutput("docker network rm ${DOCKER_NETWORK_NAME} 2> /dev/null
    commands.getoutput("docker network create --driver=bridge --attachable ${DOCKER_NETWORK_NAME} --gateway ${NETWORK_GATEWAY} --subnet ${SUBNET}

def buildContainer():
    # evalParams $@
    commands.getoutput("docker-compose -f ${composeFile} build ${buildHardOptions} ${containers}

def startContainer():
    # evalParams $@
    # source ${ENV_FILE}
    commands.getoutput("docker-compose -f ${composeFile} up -d $containers")

def showLogs():
    # evalParams $@
    # source ${ENV_FILE}
    commands.getoutput("docker-compose -f ${composeFile} logs -f $containers")


#  case switch start
def create_docker_network():
    print("hello")

def create_docker_secrets():
    print_usage()

def build():
    print("hello")

def build_hard():
    print("hello")

def run():
    print("hello")

def run_hard():
    print("hello")

def removeStacks():
    print("hello")

def removeContainers():
    createComposeFile()
    print("hello")

def show_logs_build():
    print("hello")

def show_logs_run():
    print("hello")

def quit():
    print("Goodbye")
    exit()

switcher = {
    1 : create_docker_network,
    2 : create_docker_secrets,
    3 : build,
    4 : build_hard,
    5 : run,
    6 : run_hard,
    7 : removeStacks,
    8 : removeContainers,
    9 : show_logs_build,
    10 : show_logs_run,
    11 : quit
}

def numbers_to_strings(argument):
    # Get the function from switcher dictionary
    func = switcher.get(argument, "nothing")
    # Execute the function
    return func()

def ask():
    response = None
    while response not in switcher:
        response = raw_input ("Choose a Number: ")
        response = int(response)
        if response in switcher:
            print(response)
            numbers_to_strings(response)

ask()

exit()


#
# function evalParams () {
#     containers=""
#     while (( "$#" )); do
#         if [ "$1" == "hard" ]; then
#             buildHardOptions="--force-rm --no-cache --pull"
#         elif [ "$1" == "build" ]; then
#             composeFile=${stackBuildFile}
#         elif [ "$1" == "run" ]; then
#             composeFile=${stackRunFile}
#         else
#             containers+=$1" "
#         fi
#         shift
#     done
# }
#
# function doit {
#
#         case "$1" in
#         "create_docker_network")
#             echo "create theca.cash network"
#             createDockerNetwork
#             ;;
#         "create_docker_secrets")
#             echo "create_docker_secrets"
#             ./create_docker_secrets.sh
#             ;;
#         "build")
#             echo "build theca.cash ${*:2}"
#             createComposeFile
#             buildContainer build ${*:2}
#             startContainer build ${*:2}
#             ;;
#         "build_hard")
#             echo "build theca.cash hard ${*:2}"
#             createComposeFile
#             buildContainer build hard ${*:2}
#             docker stack rm ${DOCKER_STACK_NAME_BUILD} 2> /dev/null
#             # sleeping a while to delete docker stack
#             sleep 20
# 	    startContainer build hard ${*:2}
#             ;;
#         "run")
#             echo "run theca.cash ${*:2}"
#             createComposeFile
#             buildContainer run ${*:2}
#             startContainer run ${*:2}
#             ;;
#         "run_hard")
#             echo "run theca.cash hard ${*:2}"
#             createComposeFile
#             buildContainer run hard ${*:2}
# 	    docker stack rm ${DOCKER_STACK_NAME_RUN} 2> /dev/null
#             # sleeping a while to delete docker stack
#             sleep 20
#             startContainer run hard ${*:2}
#             ;;
#         "removeContainers")
#             echo "removeContainers"
#             removeContainers
#             ;;
#         "show_logs_build")
#             echo "show logs build dockers ${*:2}"
#             showLogs build ${*:2}
#             ;;
#         "show_logs_run")
#             echo "show logs run dockers ${*:2}"
#             showLogs run ${*:2}
#             ;;
#         "Quit")
#             echo "exit"
#             exit 0
#             ;;
#         *)
#             echo "unknown option"
#             ;;
#     esac
# }
#

#
# if (( $# == 0 ))
# then
#     if (( $OPTIND == 1 ))
#     then
#         PS3='Please enter your choice: '
#         select opt in "${options[@]}"
#         do
#             doit "${opt}"
#         done
#     fi
# fi
#
# case "$1" in
#     "-b")
#         doit build ${*:2}
#         ;;
#     "-bh")
#         doit build_hard ${*:2}
#         ;;
#     "-r")
#         doit run ${*:2}
#         doit "${opt}"
#         ;;
#     "-rh")
#         doit run_hard ${*:2}
#         ;;
#     "-lb")
#         doit show_logs_build ${*:2}
#         ;;
#     "-lr")
#         doit show_logs_run ${*:2}
#         ;;
#     "-h")
#         print_usage
#         exit 0
#         ;;
#     *)
#         echo -e "argument error \n"
#         print_usage
#         exit 1
#         ;;
# esac
