#ifndef __CONFIGPARSER__
#define __CONFIGPARSER__

#include <string>
#include <fstream>
#include <streambuf>

namespace Horst {

class ConfigParser {
  protected:
    BSON::Value _config;
  public:
    ConfigParser(const std::string & filename){
        std::ifstream t(filename);
        std::string str((std::istreambuf_iterator<char>(t)),
                                std::istreambuf_iterator<char>());
        _config = BSON::Value::fromJSON(str);
    };

    void populateProcessorManager(std::shared_ptr<ProcessorManager> mgr){
        BSON::Object cfg = _config;
        for(auto & kv : cfg){
            mgr->setConfig(kv.first,kv.second);
            if(kv.second.isObject()){
                if (kv.second["class"].isString()){
                    mgr->startProcessor(kv.second["class"],kv.first);
                }
                for(auto & inner_kv : kv.second){
                    if(inner_kv.first.find("output:") == 0){
                        int outputNumber = std::stoi(inner_kv.first.substr(7));
                        if(!inner_kv.second.isString()){
                            throw std::runtime_error{"malformed config: value after output directive is not string"};
                        }
                        std::string target = inner_kv.second;
                        int colonPosition = target.find_first_of(':');
                        auto instance = target.substr(0,colonPosition);
                        int inputNumber = std::stoi(target.substr(colonPosition+1));
                        mgr->setAdjacent(kv.first,outputNumber,instance,inputNumber);
                    }
                }
            }
        }
    }

};

}

#endif // __CONFIGPARSER__