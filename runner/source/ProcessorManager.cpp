#include "horst/ProcessorManager.h"

void Horst::Processor::emit(BSON::Value && value, int output) {
    _mgr->emit(std::move(value), _id, output);
}

void Horst::ProcessorManager::startProcessor(const std::string & className, const std::string & name) {
    std::lock_guard<std::mutex> lock{_mutex};
    if(_constructors.count(className) == 0){
        std::string msg = "no such class: "+className+" (instance-name: "+name+")";
        throw std::runtime_error{msg};
    }
    _processors[name] = _constructors[className](shared_from_this(),name);
}

void Horst::ProcessorManager::stopProcessor(const std::string & name){
    std::lock_guard<std::mutex> lock{_mutex};
    if(_processors.count(name) == 0){
        throw std::runtime_error{"no such processor to stop"};
    }
    _processors.erase(name);
}


void Horst::ProcessorManager::reloadProcessor(const std::string & name, const std::string & className){
    stopProcessor(name);
    startProcessor(className, name);
}

void Horst::ProcessorManager::setAdjacent(const std::string & from, int output, const std::string & to, int input) {
    std::lock_guard<std::mutex> lock{_mutex};
    _adjacents[from + std::to_string(output)] = {to, input};
    std::cout<<"set adjacent "<<from<<":"<<output<<" to "<<to<<":"<<input<<std::endl;
}

void Horst::ProcessorManager::emit(BSON::Value && val, const std::string & from, int output) {
    std::lock_guard<std::mutex> lock{_mutex};
    auto adjacentKey = from + std::to_string(output);
    if (_adjacents.count(adjacentKey) == 0) {
        std::cout << "no adjacent for " << from <<":"<<output<< std::endl;
        return;
    }
    auto & adj = _adjacents[adjacentKey];
    auto & processor = _processors[adj.processor];
    struct Work {
        BSON::Value value;
        std::shared_ptr<Processor> processor;
        int input;
        void operator()(){
            processor->threadSafeProcess(std::move(value),input);
        };
    };
    Work work{std::move(val),processor,adj.input};
    _threadPool.add(work);
}