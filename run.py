from mininet.net import Mininet
from mininet.util import dumpNodeConnections
from mininet.log import setLogLevel

from mininet.topo import LinearTopo
from mininet.topolib import TreeTopo
from topologies.CustomFatTree import CustomFatTree

import yaml


if __name__ == '__main__':
    cfg_net = yaml.load(open('config_net.yaml', 'r'), Loader=yaml.FullLoader)
    topoName = cfg_net['topology']
    print("Mininet Topology:", topoName)

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
    print("Dumping host connections")
    dumpNodeConnections(net.hosts)
    print("Testing network connectivity")
    net.pingAll()
    net.stop()



