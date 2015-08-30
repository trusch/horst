#include "horst/ProcessorManager.h"
#include "horst/plugins/Logger.h"
#include "horst/plugins/NaturalNumbers.h"
#include "horst/plugins/Projector.h"

extern "C" void init(std::shared_ptr<Horst::ProcessorManager> mgr){

    mgr->declare("Logger",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                      const std::string & id){
        return std::shared_ptr<Horst::Processor>(
            new Horst::Logger(mgr,id)
        );
    });
    
    mgr->declare("NaturalNumbers",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                        const std::string & id){
        const auto & cfg = mgr->getConfig(id);
        int start = (long long)cfg["start"];
        int interval = (long long)cfg["interval"];
        return std::shared_ptr<Horst::Processor>(
            new Horst::NaturalNumbers(mgr,id,start,interval)
        );
    });

    mgr->declare("Projector",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                        const std::string & id){
        const auto & cfg = mgr->getConfig(id);
        BSON::Value format = cfg["format"];
        return std::shared_ptr<Horst::Processor>(
            new Horst::Projector(mgr,id,format)
        );
    });

}
