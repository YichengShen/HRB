from mininet.net import Mininet
from mininet.util import dumpNodeConnections
from mininet.log import setLogLevel

from mininet.topo import LinearTopo
from mininet.topolib import TreeTopo
from topologies.CustomFatTree import CustomFatTree

from os import system
from time import sleep
import yaml


if __name__ == '__main__':
    cfg_net = yaml.load(open('config_net.yaml', 'r'), Loader=yaml.FullLoader)
    topoName = cfg_net['topology']
    print("Mininet Topology:", topoName)

    system("sudo mn --clean")

    # Set mininet log level
    setLogLevel('info')

    # Create network
    if topoName == "LinearTopo":
        topo = LinearTopo(k=cfg_net['linear_topo']['k'],
                          n=cfg_net['linear_topo']['n'])
    elif topoName == "TreeTopo":
        topo = TreeTopo(depth=cfg_net['tree_topo']['depth'],
                        fanout=cfg_net['tree_topo']['fanout'])
    elif topoName == "CustomFatTree":
        topo = CustomFatTree(c=cfg_net['custom_fat_tree']['c'],
                             a=cfg_net['custom_fat_tree']['a'],
                             e=cfg_net['custom_fat_tree']['e'],
                             n=cfg_net['custom_fat_tree']['n'])
    net = Mininet(topo)

    # Start network
    net.start()
    # print("Dumping host connections")
    # dumpNodeConnections(net.hosts)
    # print("Testing network connectivity")
    # net.pingAll()

    hosts = topo.hosts(sort=True)
    hosts = [net.get(host) for host in hosts]
    for i, h in enumerate(hosts):
        h.cmdPrint("export PATH=$PATH:/usr/local/go/bin")
        h.cmdPrint("cd ~/go/src/HRB")
        # h.cmdPrint(f"go run main.go &")
        if i==4:
            h.cmdPrint(f"go run main.go")
        else:
            h.cmdPrint(f"go run main.go &> results/result{i}.log &")
    sleep(20)

    net.stop()



