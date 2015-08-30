#ifndef __REGEXP_FILTER__
#define __REGEXP_FILTER__

#include "horst/Processor.h"
#include <regex>

namespace Horst {

class RegexFilter : public Processor {
  protected:
    std::regex _regex;
  public:
    RegexFilter(std::shared_ptr<ProcessorManager> mgr, const std::string & id, std::string regex) : 
      Processor{mgr, id},
      _regex{regex} {}
    virtual ~RegexFilter() {}

    virtual void process(BSON::Value && value, int input) {
        if(value.isString() && std::regex_match(value.getString(),_regex)){
            emit(std::move(value));
        }
    }



};

}

#endif // __REGEXP_FILTER__