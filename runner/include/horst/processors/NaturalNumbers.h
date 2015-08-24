#ifndef __NATURALNUMBERS__
#define __NATURALNUMBERS__

#include "horst/Processor.h"
#include <thread>
#include <atomic>

namespace Horst {

class NaturalNumbersProcessor : public Processor {
  protected:
    std::thread _runloop;
    std::atomic<bool> _stop{false};
  
  public:  
    NaturalNumbersProcessor(std::shared_ptr<ProcessorManager> mgr, const std::string & id, int start, int interval) : Processor{mgr,id} {
        _runloop = std::move(std::thread{[this,start,interval](){
            int current = start;
            while(!_stop.load()){
                emit(current++);
                std::this_thread::sleep_for(std::chrono::milliseconds{interval});
            }
        }});
    }
    
    virtual ~NaturalNumbersProcessor(){
        _stop.store(true);
        if(_runloop.joinable())_runloop.join();
    }
    
    virtual void process(BSON::Value && value, int input) override {}

};

}

#endif // __NATURALNUMBERS__