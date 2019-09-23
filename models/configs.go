package models

import (
	"encoding/json"

	"github.com/appscode/go/log"
	dtypes "github.com/piersharding/r-operator/types"
	"github.com/piersharding/r-operator/utils"
	v1 "k8s.io/api/core/v1"
)

// RConfigs generates the ConfigMap for
// the R Scheduler and Worker
func RConfigs(context dtypes.RContext) (*v1.ConfigMap, error) {

	const rConfigs = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: rapp-configs-{{ .Name }}
  labels:
    app.kubernetes.io/name: rapp-configs
    app.kubernetes.io/instance: "{{ .Name }}"
    app.kubernetes.io/managed-by: MetaController
data:
  start-rapp.sh: |
    #!/usr/bin/env bash

    set -o errexit -o pipefail

    #source activate dask-distributed
    [ -f "${HOME}/.bash_profile" ] && source "${HOME}/.bash_profile"

    mkdir -p /var/log/shiny-server
    chown shiny.shiny /var/log/shiny-server
    
    if [ "$APPLICATION_LOGS_TO_STDOUT" != "false" ];
    then
        # push the "real" application logs to stdout with xtail in detached mode
        exec xtail /var/log/shiny-server/ &
    fi
    # check if the apps directory is empty - if so copy over the sample
    mkdir -p /rscripts /rlibs
    DIR_EMPTY="$(ls -A /rscripts)"
    if [ "${DIR}" == "" ]; then
        if [ -d "/srv/shiny-server/01_hello" ]; then
            cp -r /srv/shiny-server/01_hello/* /rscripts/
        fi
    fi

    # run install hook if found
    if [ -f /rscripts/rapps_install.sh ]; then
        cd /rscripts
        bash /rscripts/rapps_install.sh
    fi
    
    if [ -f /rscripts/launch.R ]; then
        cd /rscripts
        Rscript /rscripts/launch.R
    else
        #R  -e ".libPaths( c( .libPaths(), '/rlibs') ); setwd('/rscripts'); library(shiny); runApp(appDir='/rscripts', port=${PORT}, host='${HOST}', launch.browser=FALSE, display.mode='normal')" >/var/log/shiny-server/rapp.log 2>&1
        R  -e ".libPaths( c( .libPaths(), '/rlibs') ); setwd('/rscripts'); library(shiny); runApp(appDir='/rscripts', port=${PORT}, host='${HOST}', launch.browser=FALSE, display.mode='normal')" 2>&1
    fi
    
`
	result, err := utils.ApplyTemplate(rConfigs, context)
	if err != nil {
		log.Debugf("ApplyTemplate Error: %+v\n", err)
		return nil, err
	}
	configmap := &v1.ConfigMap{}
	if err := json.Unmarshal([]byte(result), configmap); err != nil {
		return nil, err
	}
	return configmap, err
}
