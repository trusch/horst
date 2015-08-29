#include "horst/ProcessorManager.h"
#include "horst/plugins/TwitterSource.h"

extern "C" void init(std::shared_ptr<Horst::ProcessorManager> mgr){
    mgr->declare("twittersource", [](std::shared_ptr<Horst::ProcessorManager> mgr, 
                                      const std::string & id){
        const auto & cfg = mgr->getConfig(id);
        auto track = cfg["track"].getString();
        auto consumerKey = cfg["consumerKey"].getString();
        auto consumerSecret = cfg["consumerSecret"].getString();
        auto accessTokenKey = cfg["accessTokenKey"].getString();
        auto accessTokenSecret = cfg["accessTokenSecret"].getString();
        return std::shared_ptr<Horst::Processor>(
            new Horst::TwitterSource(mgr,id,track,consumerKey,consumerSecret,accessTokenKey,accessTokenSecret)
        );
    });
}
