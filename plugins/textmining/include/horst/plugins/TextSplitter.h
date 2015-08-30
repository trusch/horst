#ifndef __TEXT_TO_WORDS__
#define __TEXT_TO_WORDS__

#include "horst/Processor.h"
#include <sstream>

namespace Horst {

class TextSplitter : public Processor {
  protected:
    std::string _cutset;
    bool _toLowercase;
  public:
    TextSplitter(std::shared_ptr<ProcessorManager> mgr, const std::string & id) : 
      Processor{mgr, id}{}
    virtual ~TextSplitter() {}

    virtual void process(BSON::Value && value, int input) {
        if(value.isString()){
            std::stringstream ss{value};
            std::string word;
            while(ss.good()){
              ss >> word;
              emit(word);
            }
        }
    }



};

}

#endif // __TEXT_TO_WORDS__