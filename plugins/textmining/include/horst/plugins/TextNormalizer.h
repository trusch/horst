#ifndef __TextNormalizer__
#define __TextNormalizer__

#include "horst/Processor.h"
#include <algorithm>

namespace Horst {

class TextNormalizer : public Processor {
  protected:
    std::string _cutset;
    bool _toLowercase;
  public:
    TextNormalizer(std::shared_ptr<ProcessorManager> mgr, const std::string & id, std::string cutset, bool toLowercase) : 
      Processor{mgr, id},
      _cutset{cutset},
      _toLowercase{toLowercase} {}
    virtual ~TextNormalizer() {}

    virtual void process(BSON::Value && value, int input) {
        if(value.isString()){
            std::string & str = value;
            for (unsigned int i = 0; i < _cutset.size(); ++i){
              str.erase (std::remove(str.begin(), str.end(), _cutset[i]), str.end());
            }
            if(_toLowercase){
                std::transform(str.begin(), str.end(), str.begin(), ::tolower);
            }
        }
        emit(std::move(value));
    }



};

}

#endif // __TextNormalizer__