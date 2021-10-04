from mininet.net import Mininet
from mininet.util import dumpNodeConnections
from mininet.log import setLogLevel

from mininet.topo import LinearTopo
from mininet.topolib import TreeTopo
from topologies.CustomFatTree import CustomFatTree

from os import system
from time import sleep
import yaml


def test_network(net):
    print("Waiting switch connections")
    net.waitConnected()

    print("Testing network connectivity")
    percent_loss = net.pingAll()
    assert percent_loss==0., "Network is not fully connected"

def line_count(fname):
    """Returns the number of lines in a file."""
    count = 0
    with open(fname) as f:
        for _ in f:
            count += 1
    return count


if __name__ == '__main__':
    cfg_net = yaml.load(open('config_net.yaml', 'r'), Loader=yaml.FullLoader)
    topoName = cfg_net['topology']
    print("Mininet Topology:", topoName)

    # Clean
    system("sudo mn --clean")
    system("sudo rm watch.txt")

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

    test_network(net)

    hosts = topo.hosts(sort=True)
    hosts = [net.get(host) for host in hosts]
    for i, h in enumerate(hosts):
        h.cmdPrint("export PATH=$PATH:/usr/local/go/bin")
        h.cmdPrint("cd ~/go/src/HRB")
        h.cmdPrint(f"go run main.go >> watch.txt &")

    while True:
        num_lines = line_count("watch.txt")
        print(num_lines)
        if num_lines == 5:
            break
        sleep(1)

    net.stop()


