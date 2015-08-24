#ifndef __PLUGINLOADER__
#define __PLUGINLOADER__

#include <glob.h>
#include <dlfcn.h>
#include "horst/ProcessorManager.h"


namespace Horst {

class PluginLoader {
  protected:
    std::string _directory;
    std::shared_ptr<ProcessorManager> _mgr;
  public:
    PluginLoader(const std::string & directory, std::shared_ptr<ProcessorManager> mgr) :
        _directory{directory},
        _mgr{mgr} {
        for(auto & plugin : list()){
            load(plugin);
        }
    }

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

    void load(const std::string & filename){
        void* handle = dlopen(filename.c_str(), RTLD_LAZY);
        auto init = (void (*)(std::shared_ptr<ProcessorManager>))dlsym(handle, "init");
        if (init == NULL) {
            std::string msg = "plugin failed: "+filename+": init() not found";
            throw std::runtime_error{msg};
        }
        try{
            init(_mgr);
        }catch(const std::exception & e){
            std::string msg = "plugin failed: "+filename+": init() throws error: "+e.what();
            throw std::runtime_error(msg);
        }
        std::cout<<"successfully loaded plugin "<<filename<<std::endl;
    }
};

}

#endif // __PLUGINLOADER__