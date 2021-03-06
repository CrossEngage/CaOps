# -*- mode: ruby -*-
# vi: set ft=ruby :
#
# Import current PK in ssh-agent
system('.vagrant/ssh.agent.sh')


templates = {
  'cassandra22x' => {
    dist: 'ubuntu/xenial64', mem: 1512, cpus: 1,
    copy: [ '.vagrant/default-cassandra', '.vagrant/cassandra22x.yaml', '.vagrant/CaOps.yaml', './bin/CaOps', '.vagrant/CaOps.service' ],
    exec: { '.vagrant/bootstrap.sh': [ '22x', 'C22xCluster', '10.127.22.10,10.127.22.11,10.127.22.12' ] },
  },
  'cassandra30x' => {
    dist: 'ubuntu/xenial64', mem: 1512, cpus: 1,
    copy: [ '.vagrant/default-cassandra', '.vagrant/cassandra30x.yaml', '.vagrant/CaOps.yaml', './bin/CaOps', '.vagrant/CaOps.service' ],
    exec: { '.vagrant/bootstrap.sh': [ '30x', 'C22xCluster', '10.127.30.10,10.127.30.11,10.127.30.12' ] },
  }
}

netmask = '255.255.0.0'

boxes =  {
  'c22x01' => { template: 'cassandra22x', ip: '10.127.22.11' },
  'c22x02' => { template: 'cassandra22x', ip: '10.127.22.12' },

  'c30x01' => { template: 'cassandra30x', ip: '10.127.30.11' },
  'c30x02' => { template: 'cassandra30x', ip: '10.127.30.12' },
}


Vagrant.configure(2) do |config|
  config.ssh.forward_agent = true
  config.vm.provider 'virtualbox' do |vb|
    vb.gui = false
  end

  boxes.each do |name, box|
    vmbox = templates[box[:template]].clone
    vmbox.update(box)

    config.vm.define name, autostart: true do |host|
      host.vm.hostname = name
      host.vm.box = vmbox[:dist]
      host.vm.network 'private_network', ip: vmbox[:ip], netmask: netmask

      host.vm.provider "virtualbox" do |vb|
        vb.memory = vmbox[:mem]
        vb.cpus = vmbox[:cpus]
      end

      vmbox[:copy].each do |file|
        host.vm.provision 'file', source: file, destination: '/tmp/' + File.basename(file)
      end

      vmbox[:exec].each_pair do |path, args|
        host.vm.provision 'shell', path: path.to_s, args: args + [vmbox[:ip]]
      end
    end
  end
end
