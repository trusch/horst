#ifndef __PROCESSOR__
#define __PROCESSOR__

#include <bson/Value.h>
#include <memory>

namespace Horst {

class ProcessorManager;

class Processor {
  protected:
    std::shared_ptr<ProcessorManager> _mgr;
    std::string _id;
  public:
    Processor(std::shared_ptr<ProcessorManager> mgr, const std::string & id) : _mgr{mgr}, _id{id}  {};
    void emit(BSON::Value && value, int output = 0);
    virtual void process(BSON::Value && value, int input) = 0;
    virtual ~Processor(){}
};

}

#endif // __PROCESSOR__
