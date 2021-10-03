go_tar=go1.17.1.linux-amd64.tar.gz # the version of Golang to be downloaded in install_go
mininet_ver=2.3.0                  # the version of Mininet to be downloaded in install_mininet

# Installs a version of Mininet
function install_mininet() {
    cd
    mkdir mininet
    cd mininet
    git clone git://github.com/mininet/mininet.git
    cd mininet
    git checkout ${mininet_ver}
    ./util/install.sh -a
    cd
}

# Installs a version of Golang
function install_go() {
    wget -q https://golang.org/dl/${go_tar}
    sudo tar -C /usr/local -xzf ${go_tar}
    rm ${go_tar}
    echo 'export PATH=${PATH}:/usr/local/go/bin' >>~/.bashrc
    echo 'export GOPATH=~/go' >>~/.bashrc
    echo 'export GO111MODULE="auto"' >>~/.bashrc
    source ~/.bashrc
    go version
}

install_go