#!/bin/bash
set -e

product_name=$(grep ProjectName ./config/config.go | awk -F '"' '{print $2}' | sed 's/\"//g')
Product_version_key="ProjectVersion"
VersionFile=./config/config.go
CURRENT_VERSION=$(grep ${Product_version_key} $VersionFile | awk -F '"' '{print $2}' | sed 's/\"//g')

build_path=./build
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


function build_with_type() {
    local APP_TYPE=$1

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

    build_time=$(date -u "+%Y-%m-%d_%H:%M:%S")


    local ld_flag_master="-X github.com/0xYeah/yeahBox/base_app.mGitCommitHash=${commit_hash} -X github.com/0xYeah/yeahBox/base_app.mGitCommitTime=${formatted_time} -X github.com/0xYeah/yeahBox/base_app.mGoVersion=${go_version} -X github.com/0xYeah/yeahBox/base_app.mPackageOS=${OS_TYPE} -X github.com/0xYeah/yeahBox/base_app.mPackageTime=${build_time} -X github.com/0xYeah/yeahBox/base_app.mRunMode=${RUN_MODE} -s -w"

    echo "build ${product_name}_${APP_TYPE}"



    local ld_flag_full=$ld_flag_master" -X github.com/0xYeah/yeahBox.base_app.mAppType=${APP_TYPE}"

    go build -o ${build_path}/${RUN_MODE}/${product_name}_${APP_TYPE}/${product_name}_${APP_TYPE} -trimpath -ldflags "${ld_flag_full}" ./${product_name}_${APP_TYPE}/main.go \
    && chmod a+x ${build_path}/${RUN_MODE}/${product_name}_${APP_TYPE}/${product_name}_${APP_TYPE} \
    && cp ./example_files/${product_name}_${APP_TYPE}.service ${build_path}/${RUN_MODE}/${product_name}_${APP_TYPE} \
    && cp ./example_files/install_${product_name}_${APP_TYPE}.sh ${build_path}/${RUN_MODE}/${product_name}_${APP_TYPE} \
    && mkdir -p ${build_path}/${RUN_MODE}/${product_name}_${APP_TYPE}/conf \
    && cp ./example_files/config_${APP_TYPE}.example.json ${build_path}/${RUN_MODE}/${product_name}_${APP_TYPE}/conf/config_${APP_TYPE}.json


}

function toBuild() {

    rm -rf ${build_path}/${RUN_MODE}

    mkdir -p ${build_path}/${RUN_MODE}


     echo ${build_path}/${RUN_MODE}

    build_with_type "server"
    package_files "server"

    build_with_type "agent"
    package_files "agent"
}

function package_files(){
    local APP_TYPE=$1
    local BUILD_OS_TYPE=$(echo "$OS_TYPE" | tr '[:upper:]' '[:lower:]')

    cd ${build_path}/${RUN_MODE} \
    && if [[ "$OS_TYPE" == "Windows" ]]; then
            echo "package ${product_name}_${APP_TYPE}"
            7z a ./${product_name}_${APP_TYPE}_${BUILD_OS_TYPE}_${RUN_MODE}_${CURRENT_VERSION}.zip ./${product_name}_${APP_TYPE} >/dev/null 2>&1
        else
            echo "package ${product_name}_${APP_TYPE}"
            zip -r ./${product_name}_${APP_TYPE}_${BUILD_OS_TYPE}_${RUN_MODE}_${CURRENT_VERSION}.zip ./${product_name}_${APP_TYPE}
        fi \
    && mkdir -p ../${UPLOAD_TMP_DIR} \
    && mv *.zip ../${UPLOAD_TMP_DIR} \
    && cd ../../
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

