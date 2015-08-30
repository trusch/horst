#ifndef __LOGGER_PROCESSOR__
#define __LOGGER_PROCESSOR__

#include "horst/Processor.h"
#include <chrono>
#include <iomanip>

template<typename Clock, typename Duration>
std::ostream &operator<<(std::ostream &stream,
  const std::chrono::time_point<Clock, Duration> &time_point) {
  const time_t time = Clock::to_time_t(time_point);
#if __GNUC__ > 4 || \
    ((__GNUC__ == 4) && __GNUC_MINOR__ > 8 && __GNUC_REVISION__ > 1)
  // Maybe the put_time will be implemented later?
  struct tm tm;
  localtime_r(&time, &tm);
  return stream << std::put_time(&tm, "c");
#else
  char buffer[26];
  ctime_r(&time, buffer);
  buffer[24] = '\0';  // Removes the newline that is added
  return stream << buffer;
#endif
}

namespace Horst {

class Logger : public Processor {
  protected:

  public:
    Logger(std::shared_ptr<ProcessorManager> mgr, const std::string & id) : Processor{mgr, id} {}
    virtual ~Logger() {}

    virtual void process(BSON::Value && value, int input) {
        auto now = std::chrono::system_clock::now();
        std::cout << now << " : " << _id << ":"<<input<<" > " << value.toJSON() << std::endl;
    }

};

}

#endif // __LOGGER_PROCESSOR__