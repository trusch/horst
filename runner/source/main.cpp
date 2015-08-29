#include <iostream>
#include <thread>

#include "horst/ProcessorManager.h"
#include "horst/ConfigParser.h"
#include "horst/PluginLoader.h"

using namespace Horst;

int main() {

    auto processorMgr = std::make_shared<Horst::ProcessorManager>(4);
    
    ConfigParser parser{"config.json"};
    parser.loadFile(processorMgr);
    PluginLoader pluginLoader{processorMgr->getConfig("global")["plugins"], processorMgr};
    parser.populateProcessorManager(processorMgr);

    processorMgr->join();

}