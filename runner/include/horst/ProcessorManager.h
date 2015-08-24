#ifndef __PROCESSORMANAGER__
#define __PROCESSORMANAGER__

#include <bson/Value.h>
#include <memory>
#include <mutex>

#include "horst/Processor.h"
#include "horst/ThreadPool.h"
#include "horst/ClassManager.h"
#include "horst/ConfigManager.h"

namespace Horst {

class ProcessorManager : public std::enable_shared_from_this<ProcessorManager>, public ClassManager, public ConfigManager {
    friend Processor;
  public:
    ProcessorManager(size_t numWorkers = 1) : _threadPool{numWorkers} {};
    void startProcessor(const std::string & className, const std::string & name);
    void stopProcessor(const std::string & name);
    void reloadProcessor(const std::string & name, const std::string & className);
    void setAdjacent(const std::string & from, int output, const std::string & to, int input);
    void join(){ _threadPool.join();}
    void stop(){ _threadPool.stop();}
    ~ProcessorManager(){
        stop();
        join();
    }
  protected:
    void emit(BSON::Value && val, const std::string & from, int output);

    std::mutex _mutex;
    ThreadPool _threadPool;
    std::map<std::string, std::shared_ptr<Processor>> _processors; // id -> processor
    struct Adjacent {
        std::string processor;
        int input;
    };
    std::map<std::string, Adjacent> _adjacents;
};
}

#endif // __PROCESSORMANAGER__