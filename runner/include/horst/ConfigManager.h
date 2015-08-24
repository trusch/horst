#ifndef __CONFIGMANAGER__
#define __CONFIGMANAGER__

#include <bson/Value.h>
#include <map>

namespace Horst {

class ConfigManager {
  protected:
    std::map<std::string, BSON::Value> _configs;
  public:
    const BSON::Value & getConfig(const std::string & id){
        if(_configs.count(id) == 0){
            std::string msg = "no config for processor id '"+id+"'";
            throw std::runtime_error{msg};
        }
        return _configs[id];
    };
    void setConfig(const std::string & id, BSON::Value config){
        _configs[id] = config;
    }
};

}

#endif // __CONFIGMANAGER__