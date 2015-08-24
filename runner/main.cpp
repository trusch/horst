#include <iostream>
#include <thread>

#include "horst/ProcessorManager.h"
#include "horst/ConfigParser.h"
#include "horst/processors/Logger.h"
#include "horst/processors/NaturalNumbers.h"

using namespace Horst;

int main() {

    auto processorMgr = std::make_shared<Horst::ProcessorManager>(4);
    
    processorMgr->declare("logger",[](std::shared_ptr<ProcessorManager> processorMgr, const std::string & id){
        return std::shared_ptr<Processor>(new Horst::LoggerProcessor(processorMgr,id));
    });

    processorMgr->declare("naturals",[](std::shared_ptr<ProcessorManager> processorMgr, const std::string & id){
        const auto & cfg = processorMgr->getConfig(id);
        int start = (long long)cfg["start"];
        int interval = (long long)cfg["interval"];
        return std::shared_ptr<Processor>(new NaturalNumbersProcessor(processorMgr,id,start,interval));
    });

    ConfigParser parser{"test.json"};
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