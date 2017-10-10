#!/bin/bash -v
echo "$@"

if [ "$#" -ne 4 ]; then
	echo "$0 [CASSANDRA_RELEASE] [CLUSTER_NAME] [SEEDS] [ADDRESS]"
	exit 1
fi

CASSANDRA_RELEASE="$1"
CLUSTER_NAME="$2"
SEEDS="$3"
ADDRESS="$4"

JOLOKIA_VERSION='1.3.7'

DISTRO_RELEASE=`lsb_release -s -r`
DISTRO_CODENAME=`lsb_release -s -c`

BOOTSTRAPPED_LOCK=/etc/bootstrapped.lock

if [ ! -e $BOOTSTRAPPED_LOCK ]
then
	# Disable interactive mode
	sudo mv -v /etc/apt/apt.conf.d/70debconf /root/etc-apt-apt.conf.d-70debconf.bak
	sudo dpkg-reconfigure debconf -f noninteractive -p critical

	# Disable cloud init
	echo 'datasource_list: [ None ]' | sudo -s tee /etc/cloud/cloud.cfg.d/90_dpkg.cfg
	sudo dpkg-reconfigure -f noninteractive cloud-init

	# Uninstall unused software
	sudo apt-get purge chef chef-zero puppet puppet-common landscape-client landscape-common -y

	# Repos
	echo "deb http://www.apache.org/dist/cassandra/debian $CASSANDRA_RELEASE main" | \
		sudo tee /etc/apt/sources.list.d/cassandra.sources.list
	wget -O - https://www.apache.org/dist/cassandra/KEYS | sudo apt-key add -

	sudo apt-get update

	# Upgrade everything
	sudo apt-get upgrade -y && sudo apt-get dist-upgrade -y && sudo apt-get autoremove -y

	# Install base pkgs
	sudo apt-get install language-pack-en -y

	# Setup Timezone
	echo "Europe/Berlin" | sudo tee /etc/timezone
	sudo dpkg-reconfigure --frontend noninteractive tzdata

	# Java
	sudo apt-get -y install openjdk-8-jdk

	# Jolokia
	wget "https://repo1.maven.org/maven2/org/jolokia/jolokia-jvm/$JOLOKIA_VERSION/jolokia-jvm-$JOLOKIA_VERSION-agent.jar" \
		-O /usr/share/java/jolokia-jvm-agent.jar

	# Cassandra
	sudo RUNLEVEL=1 DEBIAN_FRONTEND=noninteractive DEBIAN_PRIORITY=critical \
		apt-get -q -y -o "Dpkg::Options::=--force-confdef" -o "Dpkg::Options::=--force-confold" install \
			cassandra cassandra-tools

	# CaOps
	sudo mkdir /etc/CaOps
fi

diff -qs /tmp/default-cassandra /etc/default/cassandra
DIFF_CASS_DEFAULT="$?"
cat /tmp/default-cassandra | sudo tee /etc/default/cassandra

cat /tmp/cassandra${CASSANDRA_RELEASE}.yaml | \
	sed "s/{{CLUSTER_NAME}}/$CLUSTER_NAME/g" | \
	sed "s/{{SEEDS}}/$SEEDS/g" | \
	sed "s/{{LISTEN_ADDRESS}}/$ADDRESS/g" | \
	sed "s/{{RPC_ADDRESS}}/$ADDRESS/g" | \
	tee /tmp/cassandra.yaml
diff -qs /tmp/cassandra.yaml /etc/cassandra/cassandra.yaml
DIFF_CASS_CONF="$?"
cat /tmp/cassandra.yaml | sudo tee /etc/cassandra/cassandra.yaml

if [ "$DIFF_CASS_DEFAULT" == "1" ] || [ "$DIFF_CASS_CONF" == "1" ]; then
	sudo systemctl restart cassandra
fi

cat /tmp/CaOps.yaml | \
	sed "s/{{IP}}/$ADDRESS/g" | \
	sudo tee /etc/CaOps/CaOps.yaml

cat /tmp/CaOps.service | sudo tee /etc/systemd/system/CaOps.service
sudo systemctl daemon-reload
systemctl enable CaOps.service

sudo systemctl stop CaOps.service
sudo cp /tmp/CaOps /usr/local/bin/CaOps
sudo systemctl start CaOps.service

if [ ! -e $BOOTSTRAPPED_LOCK ]
then
	sudo touch $BOOTSTRAPPED_LOCK
	echo "Going to reboot in 10 seconds..."
	sleep 10
	sudo reboot
fi
