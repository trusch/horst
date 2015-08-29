#include "horst/ProcessorManager.h"
#include "horst/plugins/Logger.h"
#include "horst/plugins/NaturalNumbers.h"
#include "horst/plugins/Projector.h"

extern "C" void init(std::shared_ptr<Horst::ProcessorManager> mgr){
    mgr->declare("logger",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                      const std::string & id){
        return std::shared_ptr<Horst::Processor>(
            new Horst::LoggerProcessor(mgr,id)
        );
    });
    
    mgr->declare("naturals",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                        const std::string & id){
        const auto & cfg = mgr->getConfig(id);
        int start = (long long)cfg["start"];
        int interval = (long long)cfg["interval"];
        return std::shared_ptr<Horst::Processor>(
            new Horst::NaturalNumbersProcessor(mgr,id,start,interval)
        );
    });

    mgr->declare("projector",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                        const std::string & id){
        const auto & cfg = mgr->getConfig(id);
        BSON::Value format = cfg["format"];
        std::cout<<"construct projector"<<std::endl;
        return std::shared_ptr<Horst::Processor>(
            new Horst::ProjectorProcessor(mgr,id,format)
        );
    });
}
