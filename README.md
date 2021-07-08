# Lemur with Azure Blob Storage support

## Introduction

This project is a version of Lemur providing a copy tool that supports Azure Blob Storage. 

Lustre has a concept called Hierarchical Storage Management, or HSM. HSM is a backend to Lustre, that could for example be a POSIX file store, a tape system, or modern blob storage. For Lustre to be able to move data between Lustre and HSM it needs a copy tool. Lemur is such a tool.

## How to build the binaries without Docker

### Prerequisites

In order to build the binaries for Lemur you need to have the following available:

* A machine with the `lustre-client` package installed, CentOS 7.x or RHEL 7.x work well. To install from rpm use `yum install -y lustre-client`.
* A recent version of Go. If you need a quick way to install go, look at [canha/golang-tools-install-script](https://github.com/canha/golang-tools-install-script).
* An active GOPATH set up

If you don't want to do the prep for the environment, then consider using the [Terraform templates](https://github.com/Azure/azlustre/tree/main/terraform) that build a Lustre cluster for you.

### Build instructions

This repository has been tested with Lustre 2.12 and higher, but will compile for 2.10 too if needed. If you need 2.10, please adjust the CGO_CFLAGS to point to the correct 2.10 headers.

```bash
git clone https://github.com/wastore/lemur.git
cd lemur/
export CGO_CFLAGS='-I/usr/src/lustre-2.12.5/lustre/include -I/usr/src/lustre-2.12.5/lustre/include/uapi'
go build -v -i -ldflags "-X 'main.version=0.6.0_56_g286df59_dirty'" -o dist/lhsm-plugin-az ./cmd/lhsm-plugin-az
```

This will create a binary `lhsm-plugin-az` in the `lemur/dist/` directory.



## How to use this tool

### Prerequisites 

* Set up a Lustre system, including MGS, MDT and OSS nodes.
* You will need a separate node for the copy tool, because the copy tool basically functions as a client of Lustre. Install a lustre-client on this node.

### Configuration

#### Ensure HSM is enabled, it is not enabled by default

You need to enable HSM on your Lustre cluster. Depending on your set up this may not be enabled yet. To enable the HSM please execute the following bits from one of your Lustre client machines (e.g. your target HSM machine) where your Lustre file system is currently mounted.

```bash
lctl set_param -P mdt.*-MDT0000.hsm_control=enabled
lctl set_param -P mdt.*-MDT0000.hsm.default_archive_id=1
lctl set_param mdt.*-MDT0000.hsm.max_requests=128
```

#### Configure the plugin

Copy your `lhsm-plugin-az` binary to a plugin directory, typically `/usr/libexec/lhsmd`. You will need to set up the lhsmd agent at `/etc/lhsmd/agent` and the configuration for the azure plugin at `/etc/lhsmd/lhsm-plugin-az`. 

A sample `/etc/lhsmd/agent` is:

```
# Lustre NID and filesystem name for the front end filesystem, the agent will mount this.
# This is <mgs_ip>@tcp:/<fs_name>
client_device="10.10.4.6@tcp:/lustrefs" 

enabled_plugins=["lhsm-plugin-az"] 

# Directory to look for the plugins
plugin_dir="/usr/libexec/lhsmd"

handler_count=16

snapshots {
        enabled = false
}
```

A sample `/etc/lhsmd/lhsm-plugin-az` looks like this:
```
az_storage_account = "$storage_account"
az_storage_key = "$storage_key"
num_threads = 32

# One or more archive definition is required.

archive "az-blob" {
    id = 1                           # Must be unique to this endpoint
    container = "$storage_container" # Container used for this archive
    prefix = ""                      # Optional prefix
    num_threads = 32
}
```

#### Validate that it all works

Fire up lhsmd using the `--debug` flag and check for any errors. 

#### Set up LHSMD as a service

You can run LHSMD from your shell directly, in a terminal emulator, but it is easier to do as a service. Set LHSMD as a service with a systemd script and start the service.


```bash
cat <<EOF >/etc/systemd/system/lhsmd.service
[Unit]
Description=The lhsmd server
After=syslog.target network.target remote-fs.target nss-lookup.target
[Service]
Type=simple
PIDFile=/run/lhsmd.pid
ExecStartPre=/bin/mkdir -p /var/run/lhsmd
ExecStart=/sbin/lhsmd -config /etc/lhsmd/agent
Restart=always
[Install]
WantedBy=multi-user.target
EOF

chmod 600 /etc/systemd/system/lhsmd.service

systemctl daemon-reload
systemctl enable lhsmd
systemctl start lhsmd
```

### Usage

Once installed you can use `lfs hsm_*` commands to interact with HSM and move a file to Azure Blob Storage and back. 

## How to build executable for CBLD environment and run tests

This doc briefly goes over the steps involved in building the plugin for testing inside a deployed Lustre cluster.

### Preparation

- Install Docker
- Run the following commands:
  - az login
  - az acr login --name [container-registry-name]
    - This is to allow us to pull images from a private registry called `[container-registry-name]`.
  - docker pull [container-registry-name].azurecr.io/copytoolbuildimage/gobuild

### To build the executable

- docker run --name gobuild -v /{path-to-your-lemur-project-code}:/usr/src/lemur -it [container-registry-name].azurecr.io/copytoolbuildimage/gobuild:latest
  - Explanation:
    - We are starting a container that has all the right dependencies built (especially Lustre itself)
    - The -v flag is mapping a path on your machine (the lemur project) into a specific path on the container, so that the source code shows up inside the container.
    - Any changes you make in /usr/src/lemur on the laptop OR in the container should persist.
- cd /usr/src/lemur
  - You should see your lemur project's source code now
- ./build_plugin.sh
  - This step should run for a while, and the output will be in the `dist` folder

### To clean up 

- Break out of the container with ‘exit’.
- docker rm gobuild
  - This gets rid of the container so that you can `docker run` the next time, otherwise the name will be occupied.
 
### How to debug if something goes wrong with compilation

- The environment variables are set appropriately already, but you can examine them with ‘set’ command.
- The most important dependencies to check are:
  - The CGO_CFLAGS which specifies where to look for the C header files
  - The Lustre binaries themselves have to be available to the C compiler, check to make sure
- Check with Joe to see if the image has changed in unexpected ways.

### How to test the compiled executable

- Deploy a cluster as the guide specifies
- SSH into the HSM VM (the name starts with 'APRI') by looking up its private IP address in the portal
- Transfer your newly compiled build to that machine (you can do so via a storage container)
- See if the current lhsmd is running:
  - ps -A | grep lhsmd
- Stop the current lhsmd:
  - sudo systemctl stop lhsmd
- Move your new build into /usr/laaso/bin/ with a new name
- Edit the agent config: sudo vi /etc/lhsmd/agent
  - point it to use the new build instead
- Add a new config for the new build: cp /etc/lhsmd/lhsm-plugin-az /etc/lhsmd/lhsm-plugin-new-name
- Either:
  - restart the lhsmd with: sudo systemctl start lhsmd
  - Or better yet run it directly in a separate window: /usr/laaso/bin/lhsmd -config /etc/lhsmd/agent
    - This allows you to see the output directly
- cd /lustre/client
- Exercise the copy tool by triggering restore operations 


### Extra notes

- The lhsmd config is located at: sudo vi /etc/systemd/system/lhsmd.service
- The plugin config is located at: sudo vi /etc/lhsmd/lhsm-plugin-az
  - The name of the config is the same as the plugin.
  - If you've added new fields to the plugin config, you need to update this file before running lhsmd.
- lhsmd is the parent that starts the plugin, so stopping lhsmd also kills the copy tool
- The sys logs are located at /var/log/daemon.log
- You can trigger multiple file restores:
  - Example:
    ```bash
    IFS=$'\n'; set -f
    for f in $(find /lustre/ST0202 -type f); do lfs hsm_restore "$f"; done
    unset IFS; set +f
    ```
 
