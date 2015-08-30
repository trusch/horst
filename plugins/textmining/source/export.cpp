#include "horst/ProcessorManager.h"
#include "horst/plugins/TextNormalizer.h"
#include "horst/plugins/TextSplitter.h"
#include "horst/plugins/RegexFilter.h"
#include "horst/plugins/IncrementalWordHistogram.h"

extern "C" void init(std::shared_ptr<Horst::ProcessorManager> mgr){
    
    mgr->declare("TextNormalizer",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                        const std::string & id){
        const auto & cfg = mgr->getConfig(id);
        BSON::Value cutset = cfg["cutset"];
        BSON::Value tolower = cfg["tolower"];
        if(!cutset.isString()){
            cutset.reset();
            cutset = "!?.,\"'";
        }
        if(!tolower.isBool()){
            tolower.reset();
            tolower = false;
        }
        return std::shared_ptr<Horst::Processor>(
            new Horst::TextNormalizer(mgr,id,cutset,tolower)
        );
    });

    mgr->declare("TextSplitter",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                        const std::string & id){
        return std::shared_ptr<Horst::Processor>(
            new Horst::TextSplitter(mgr,id)
        );
    });

    mgr->declare("RegexFilter",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                        const std::string & id){
        const auto & cfg = mgr->getConfig(id);
        std::string regex = cfg["regex"];
        return std::shared_ptr<Horst::Processor>(
            new Horst::RegexFilter(mgr,id,regex)
        );
    });

    mgr->declare("IncrementalWordHistogram",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                        const std::string & id){
        const auto & cfg = mgr->getConfig(id);
        double decayFactor = cfg["decayFactor"];
        double lowerThreshold = cfg["lowerThreshold"];
        return std::shared_ptr<Horst::Processor>(
            new Horst::IncrementalWordHistogram(mgr,id,decayFactor,lowerThreshold)
        );
    });

}
