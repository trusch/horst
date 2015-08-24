#include <iostream>
#include <thread>

#include "horst/ProcessorManager.h"
#include "horst/ConfigParser.h"
#include "horst/PluginLoader.h"

using namespace Horst;

int main() {

    auto processorMgr = std::make_shared<Horst::ProcessorManager>(4);
    
    ConfigParser parser{"test.json"};
    parser.loadFile(processorMgr);
    PluginLoader pluginLoader{processorMgr->getConfig("global")["plugins"], processorMgr};
    parser.populateProcessorManager(processorMgr);

    while(true){
        std::this_thread::sleep_for(std::chrono::milliseconds{1000});
        processorMgr->setAdjacent("my_nats_1",0,"my_logger_2",0);
        processorMgr->setAdjacent("my_nats_2",0,"my_logger_2",0);
        std::this_thread::sleep_for(std::chrono::milliseconds{1000});
        processorMgr->setAdjacent("my_nats_1",0,"my_logger_1",0);
        processorMgr->setAdjacent("my_nats_2",0,"my_logger_1",0);
    }

    processorMgr->join();

}