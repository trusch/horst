#ifndef __PROCESSOR__
#define __PROCESSOR__

#include <bson/Value.h>
#include <memory>
#include <mutex>

namespace Horst {

class ProcessorManager;

class Processor {
  friend ProcessorManager;
  protected:
    std::shared_ptr<ProcessorManager> _mgr;
    std::string _id;
    std::mutex _mutex;
    void threadSafeProcess(BSON::Value && value, int input){
        std::lock_guard<std::mutex> lock{_mutex};
        process(std::move(value), input);
    }
  public:
    Processor(std::shared_ptr<ProcessorManager> mgr, const std::string & id) : _mgr{mgr}, _id{id}  {};
    void emit(BSON::Value && value, int output = 0);
    virtual void process(BSON::Value && value, int input) = 0;
    virtual ~Processor(){}
};

}

#endif // __PROCESSOR__
