#ifndef __INCREMENTAL_WORD_HISTOGRAM__
#define __INCREMENTAL_WORD_HISTOGRAM__

#include "horst/Processor.h"
#include <map>

namespace Horst {

class IncrementalWordHistogram : public Processor {
  protected:
    double _decayFactor;
    double _lowerThreshold;
    std::map<std::string,double> _histogram;
    double _sum;
    void update(std::string & word){
      double & value = _histogram[word];
      value += 1.0;
      _sum += 1.0;
    }
    void decay(){
      std::vector<std::string> toDelete;
      _sum *= _decayFactor;
      for(auto & kv : _histogram){
        _histogram[kv.first] = kv.second * _decayFactor;
        if(kv.second < _lowerThreshold){
          toDelete.push_back(kv.first);
        }
      }
      for(auto & key : toDelete){
        _histogram.erase(key);
      }
    }
    void emitHistogram(){
      BSON::Value val;
      for(auto & kv : _histogram){
        val[kv.first] = (kv.second/_sum)*100;
      }
      emit(std::move(val));
    }
  public:
    IncrementalWordHistogram(std::shared_ptr<ProcessorManager> mgr, const std::string & id, double decayFactor, double lowerThreshold) : 
      Processor{mgr, id},
      _decayFactor{decayFactor},
      _lowerThreshold{lowerThreshold} {}
    virtual ~IncrementalWordHistogram() {}

    virtual void process(BSON::Value && value, int input) {
        if(value.isString()){
            update(value);
            decay();
            emitHistogram();
        }
    }



};

}

#endif // __INCREMENTAL_WORD_HISTOGRAM__