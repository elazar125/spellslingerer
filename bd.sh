#! /usr/bin/bash

base_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

binary=""
migrations=""
templates=""
javascript=""
css=""
vendor=""

while getopts "btmjcv" OPTION
do
    case "$OPTION" in
        b)
            binary="true"
            ;;
        t)
            templates="true"
            ;;
        m)
            migrations="true"
            ;;
        j)
            javascript="true"
            ;;
        c)
            css="true"
            ;;
        v)
            vendor="true"
            ;;
        ?)
            echo "
Build and Deploy Spellslingerer
Specify portions to deploy with flags:
  -b: binary
  -t: templates
  -m: migrations
  -j: javascript
  -c: css
  -v: vendor
"
            exit 1
            ;;
    esac
done


build_deploy_binary() {
  env CGO_ENABLED=0 GOOS=linux go build -o dist/spellslingerer.new
  scp dist/spellslingerer.new spellslingerer:~/spellslingerer
  ssh spellslingerer -t "
    systemctl stop spellslingerer
    mv ~/spellslingerer/spellslingerer.new ~/spellslingerer/spellslingerer
    chmod +x ~/spellslingerer/spellslingerer
    systemctl start spellslingerer
  "
}

deploy_templates() {
  scp -r templates spellslingerer:~/spellslingerer
  ssh spellslingerer -t "systemctl restart spellslingerer"
}

deploy_migrations() {
  scp -r migration_files spellslingerer:~/spellslingerer
}

deploy_javascript() {
  scp -r pb_public/js spellslingerer:~/spellslingerer/pb_public
}

deploy_css() {
  scp -r pb_public/css spellslingerer:~/spellslingerer/pb_public
}

deploy_vendor() {
  scp -r pb_public/vendor spellslingerer:~/spellslingerer/pb_public
}

pushd "$base_dir" || exit

if [[ -n "$javascript" ]]
then
  deploy_javascript
fi

if [[ -n "$vendor" ]]
then
  deploy_vendor
fi

if [[ -n "$css" ]]
then
  deploy_css
fi

if [[ -n "$templates" ]]
then
  deploy_templates
fi

if [[ -n "$migrations" ]]
then
  deploy_migrations
fi

if [[ -n "$binary" ]]
then
  build_deploy_binary
fi

popd || exit
