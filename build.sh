#!/bin/bash
set -e

product_name=$(grep ProjectName ./config/config.go | awk -F '"' '{print $2}' | sed 's/\"//g')
Product_version_key="ProjectVersion"
VersionFile=./config/config.go
CURRENT_VERSION=$(grep ${Product_version_key} $VersionFile | awk -F '"' '{print $2}' | sed 's/\"//g')

build_path=../build
RUN_MODE=release

UPLOAD_TMP_DIR=upload_tmp_dir

OS_TYPE="Unknown"
GetOSType(){
    uNames=`uname -s`
    osName=${uNames: 0: 4}
    if [ "$osName" == "Darw" ] # Darwin
    then
        OS_TYPE="Darwin"
    elif [ "$osName" == "Linu" ] # Linux
    then
        OS_TYPE="Linux"
    elif [ "$osName" == "MING" ] # MINGW, windows, git-bash
    then
        OS_TYPE="Windows"
    else
        OS_TYPE="Unknown"
    fi
}
GetOSType

function toBuild() {

    rm -rf ${build_path}/${RUN_MODE}
    mkdir -p ${build_path}/${RUN_MODE}

    go_version=$(go version | awk '{print $3}')
    commit_hash=$(git show -s --format=%H)
    commit_date=$(git show -s --format="%ci")

    if [[ "$OS_TYPE" == "Darwin" ]]; then
        # macOS
        formatted_time=$(date -u -j -f "%Y-%m-%d %H:%M:%S %z" "${commit_date}" "+%Y-%m-%d_%H:%M:%S")
    else
        # Linux
        formatted_time=$(date -u -d "${commit_date}" "+%Y-%m-%d_%H:%M:%S")
    fi

    build_time=$(date +"%Y-%m-%d_%H:%M:%S")

    local ld_flag_master="-X main.mGitCommitHash=${commit_hash} -X main.mGitCommitTime=${formatted_time} -X main.mGoVersion=${go_version} -X main.mPackageOS=${OS_TYPE} -X main.mPackageTime=${build_time} -X main.mRunMode=${RUN_MODE} -s -w"

    cd ./pre_app \
    && go mod tidy \
    && cd ../

    echo "build ${product_name}_server"
    local ld_flag_server=$ld_flag_master" -X main.mRunWith=server"

    cd ./pre_app \
    && go build -o ${build_path}/${RUN_MODE}/${product_name}_server/${product_name}_server -trimpath -ldflags "${ld_flag_server}" ./main.go \
    && chmod a+x ${build_path}/${RUN_MODE}/${product_name}_server/${product_name}_server \
    && cp ./example_files/${product_name}_server.service ${build_path}/${RUN_MODE}/${product_name}_server \
    && cp ./example_files/install_${product_name}_server.sh ${build_path}/${RUN_MODE}/${product_name}_server \
    && mkdir -p ${build_path}/${RUN_MODE}/${product_name}_server/conf \
    && cp ./example_files/config_server.example.json ${build_path}/${RUN_MODE}/${product_name}_server/conf/config_server.json \
    && cd ../


    echo "build ${product_name}_agent"
    local ld_flag_agent=$ld_flag_master" -X main.mRunWith=agent"

    cd ./pre_app \
    && go build -o ${build_path}/${RUN_MODE}/${product_name}_agent/${product_name}_agent -trimpath -ldflags "${ld_flag_agent}" ./main.go \
    && chmod a+x ${build_path}/${RUN_MODE}/${product_name}_agent/${product_name}_agent \
    && cp ./example_files/${product_name}_agent.service ${build_path}/${RUN_MODE}/${product_name}_agent \
    && cp ./example_files/install_${product_name}_agent.sh ${build_path}/${RUN_MODE}/${product_name}_agent \
    && mkdir -p ${build_path}/${RUN_MODE}/${product_name}_agent/conf \
    && cp ./example_files/config_agent.example.json ${build_path}/${RUN_MODE}/${product_name}_agent/conf/config_agent.json


#    package_files
}

function package_files(){

    local BUILD_OS_TYPE=$(echo "$OS_TYPE" | tr '[:upper:]' '[:lower:]')

    cd ${build_path}/${RUN_MODE} \
    && if [[ "$OS_TYPE" == "Windows" ]]; then
            echo "package ${product_name}_server"
            7z a ./${product_name}_server_${BUILD_OS_TYPE}_${RUN_MODE}_${CURRENT_VERSION}.zip ./${product_name}_server >/dev/null 2>&1
            echo "package ${product_name}_agent"
            7z a ./${product_name}_agent_${BUILD_OS_TYPE}_${RUN_MODE}_${CURRENT_VERSION}.zip ./${product_name}_agent >/dev/null 2>&1
        else
            echo "package ${product_name}_server"
            zip -r ./${product_name}_server_${BUILD_OS_TYPE}_${RUN_MODE}_${CURRENT_VERSION}.zip ./${product_name}_server
            echo "package ${product_name}_agent"
            zip -r ./${product_name}_agent_${BUILD_OS_TYPE}_${RUN_MODE}_${CURRENT_VERSION}.zip ./${product_name}_agent
        fi \
    && mkdir -p ../${UPLOAD_TMP_DIR} \
    && mv *.zip ../${UPLOAD_TMP_DIR} \
    && cd ../
}


function handlerunMode() {
    if [[ "$1" == "release" || "$1" == "" ]]; then
        RUN_MODE=release
    elif [[ "$1" == "test" ]]; then
        RUN_MODE=test
    elif [[ "$1" == "debug" ]]; then
        RUN_MODE=debug
    else
        echo "Usage: bash build.sh [release|test],default with:release"
        exit 0
    fi
}


handlerunMode "$1" && toBuild

