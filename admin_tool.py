#!/usr/bin/python
# -*- coding: utf-8 -*-

import subprocess
import time

ENV_FILE = ".env"
# source ${ENV_FILE}

stackCommonServicesFile = "docker-compose.stack_common-services.yml"
stackBuildFile = "docker-compose.stack_build.yml"
stackRunFile = "docker-compose.stack_run.yml"

def print_usage():
    print("Usage: admin_tool.sh [option] [OPTIONAL_DOCKER_TO_BUILD (default: all) docker-service]\n"
            "option:  -b        to build containers and build theca.cash in docker\n"
            "         -bh       to build containers and build theca.cash in docker (build --force-rm --no-cache --pull)\n"
            "         -r        to build containers and run theca.cash in docker\n"
            "         -rh       to build containers and run theca.cash in docker (build --force-rm --no-cache --pull) \n"
            "         -lb       to docker service logs from build dockers\n"
            "         -lr       to docker service logs from run dockers\n")

def createComposeFile():
    buildSecretConfig = "docker-compose.build_template.yml"
    runSecretConfig = "docker-compose.run_template.yml"
    print(subprocess.call("docker-compose -f docker-compose.build_template.yml -f %s config > %s" % (buildSecretConfig, stackBuildFile), shell=True))
    print(subprocess.call("docker-compose -f docker-compose.run_template.yml -f %s config > %s" % (runSecretConfig, stackRunFile), shell=True))

def evalParams(para):
    containers = ""
    buildHardOptions = ""
    composeFile = ""
    if para == "hard":
        buildHardOptions = "--force-rm --no-cache --pull"
    elif para == "build":
        composeFile = stackBuildFile
    elif para == "run":
        composeFile = stackRunFile
    # else:
        # containers = $1" "
    return composeFile, buildHardOptions, containers

# def removeStacks():
    # print(subprocess.call("docker stack rm %s" % DOCKER_STACK_NAME_BUILD))
    # print(subprocess.call("docker stack rm %s" % DOCKER_STACK_NAME_RUN))

def removeContainers():
    print(subprocess.call("createComposeFile compose"))
    print(subprocess.call("docker-compose -f %s -f %s rm -sf" % (stackBuildFile, stackRunFile)))

# def createDockerNetwork():
#     subprocess.call("docker network rm ${DOCKER_NETWORK_NAME} 2> /dev/null")
#     subprocess.call("docker network create --driver=bridge --attachable ${DOCKER_NETWORK_NAME} --gateway ${NETWORK_GATEWAY} --subnet ${SUBNET}")

def buildContainer(do):
    composeFile, buildHardOptions, containers = evalParams(do)
    print(subprocess.call("docker-compose -f %s build %s %s}" % (composeFile, buildHardOptions, containers), shell=True))

def startContainer(do):
    composeFile, buildHardOptions, containers  = evalParams(do)
    print(subprocess.call("docker-compose -f %s up -d $containers" % composeFile, shell=True))

def showLogs(do):
    composeFile, buildHardOptions, containers = evalParams(do)
    print(subprocess.call("docker-compose -f %s logs -f $containers" % composeFile, shell=True))

def create_docker_secrets():
    print("create_docker_secrets")
    print(subprocess.call("./create_docker_secrets.sh"))

def build():
    print("build theca.cash")
    createComposeFile()
    buildContainer("build")
    startContainer("build")

def build_hard():
    print("build theca.cash hard")
    createComposeFile()
    buildContainer("build")
    # print(subprocess.call("docker stack rm % 2> /dev/null" % DOCKER_STACK_NAME_BUILD))
    # sleeping a while to delete docker stack
    time.sleep(20)

def run():
    print("run theca.cash")
    createComposeFile()
    buildContainer("run")
    startContainer("run")

def run_hard():
    print("run theca.cash hard")
    createComposeFile()
    buildContainer("run hard")
    # removeStacks()

def removeStacks():
    # print(subprocess.call("docker stack rm %s 2> /dev/null" % DOCKER_STACK_NAME_RUN))
    # sleeping a while to delete docker stack
    time.sleep(20)
    startContainer("run hard")

def removeContainers():
    createComposeFile()
    print("removeContainers")
    removeContainers()

def show_logs_build():
    print("show logs build dockers")
    showLogs("build")

def show_logs_run():
    print("show logs run dockers")
    showLogs("run")

def quit():
    print("Goodbye")
    exit()

switcher = {
    # 1 : create_docker_network,
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
        response = input ("Choose a Number: ")
        response = int(response)
        if response in switcher:
            # print(response)
            numbers_to_strings(response)

# print_usage()

ask()
