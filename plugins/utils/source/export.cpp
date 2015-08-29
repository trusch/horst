#include "horst/ProcessorManager.h"
#include "horst/plugins/Logger.h"
#include "horst/plugins/NaturalNumbers.h"
#include "horst/plugins/Projector.h"
#include "horst/plugins/StringNormalizer.h"
#include "horst/plugins/TextToWords.h"
#include "horst/plugins/RegexFilter.h"
#include "horst/plugins/IncrementalWordHistogram.h"

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
        return std::shared_ptr<Horst::Processor>(
            new Horst::ProjectorProcessor(mgr,id,format)
        );
    });

    mgr->declare("stringnormalizer",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
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
            new Horst::StringNormalizer(mgr,id,cutset,tolower)
        );
    });

    mgr->declare("texttowords",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                        const std::string & id){
        return std::shared_ptr<Horst::Processor>(
            new Horst::TextToWordsProcessor(mgr,id)
        );
    });

    mgr->declare("regexfilter",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                        const std::string & id){
        const auto & cfg = mgr->getConfig(id);
        std::string regex = cfg["regex"];
        return std::shared_ptr<Horst::Processor>(
            new Horst::RegexFilter(mgr,id,regex)
        );
    });

    mgr->declare("incrementalwordhistogram",[](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                        const std::string & id){
        const auto & cfg = mgr->getConfig(id);
        double decayFactor = cfg["decayFactor"];
        double lowerThreshold = cfg["lowerThreshold"];
        return std::shared_ptr<Horst::Processor>(
            new Horst::IncrementalWordHistogram(mgr,id,decayFactor,lowerThreshold)
        );
    });

    
}
