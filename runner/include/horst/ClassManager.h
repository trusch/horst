#ifndef __CLASSMANAGER__
#define __CLASSMANAGER__

#include <map>
#include <functional>
#include "horst/ProcessorManager.h"

namespace Horst {

typedef std::function<std::shared_ptr<Processor>(std::shared_ptr<ProcessorManager> mgr, const std::string & id)> ProcessorConstructor;

class ClassManager {
  protected:
    std::map<std::string, ProcessorConstructor> _constructors;
  public:
    void declare(const std::string & className, ProcessorConstructor constructor){
        _constructors[className] = constructor;
    };
    std::shared_ptr<Processor> construct(const std::string & className, std::shared_ptr<ProcessorManager> mgr, const std::string & id){
        if(_constructors.count(className) == 0){
            throw std::runtime_error{"no such processor class"};
        }
        auto & constructor = _constructors[className]; 
        return constructor(mgr, id);
    }
};

}

#endif // __CLASSMANAGER__