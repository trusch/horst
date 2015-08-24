#ifndef __PLUGINLOADER__
#define __PLUGINLOADER__

#include <glob.h>

namespace Horst {

class PluginLoader {
  protected:
    std::string _directory;
  public:
    PluginLoader(const std::string & directory) : _directory{directory} {}

    std::vector<std::string> list(){
        std::vector<std::string> result;
        glob_t globbuf;
        int err = glob((_directory+"/*.so").c_str(), 0, NULL, &globbuf);
        if(err == 0) {
            for (size_t i = 0; i < globbuf.gl_pathc; i++) {
                result.push_back(globbuf.gl_pathv[i]);
            }
            globfree(&globbuf);
        }
        return result;
    }

    bool load(const std::string & filename){
        
    }

};

}

#endif // __PLUGINLOADER__