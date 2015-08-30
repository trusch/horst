#ifndef __PROJECTOR_PROCESSOR__
#define __PROJECTOR_PROCESSOR__

#include "horst/Processor.h"

namespace Horst {

class Projector : public Processor {
  protected:
    BSON::Value _format;

    BSON::Value getByDotNotation(BSON::Value & val, std::string key){
        std::deque<std::string> keys;
        std::stringstream ss(key);
        std::string item;
        while (std::getline(ss, item, '.')) {
            keys.push_back(item);
        }
        return getByDotNotation(val, keys);
    }
    
    BSON::Value getByDotNotation(BSON::Value & val, std::deque<std::string> & keys){
        if(keys.size() == 0){
            return val;
        }
        if(val.isObject()){
            std::string key = keys.front();
            keys.pop_front();
            return getByDotNotation(val[key], keys);
        }
        return BSON::Value{};
    }

    BSON::Value fill(BSON::Value & val, BSON::Value & src){
        if(val.isString()){
            std::string & str = val;
            if(str[0] == '@'){
                return getByDotNotation(src, str.substr(1));
            }
        }
        if(val.isObject()){
            for(auto & kv : val.getObject()){
                val[kv.first] = fill(kv.second,src);
            }
        }
        return val;
    }

  public:
    Projector(std::shared_ptr<ProcessorManager> mgr, const std::string & id, BSON::Value format) : 
      Processor{mgr, id},
      _format{format} {}
    virtual ~Projector() {}

    virtual void process(BSON::Value && value, int input) {
        BSON::Value result{_format};
        result = fill(result,value);
        emit(std::move(result));
    }



};

}

#endif // __PROJECTOR_PROCESSOR__