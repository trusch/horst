#ifndef __TWITTERSOURCE__
#define __TWITTERSOURCE__

#include "horst/Processor.h"
#include "horst/plugins/LineFramer.h"
#include <curl/curl.h>
#include <oauth.h>

namespace Horst {

class TwitterSource : public Processor {
  protected:
    struct Proc
    {
        TwitterSource *_twitterSource;
        const char* cUrl;
        const char* cConsKey;
        const char* cConsSec;
        const char* cAtokKey;
        const char* cAtokSec;
        CURL        *curl;
        char*       cSignedUrl;
        Proc(TwitterSource*, const char*, const char*, const char*, const char*);
        void execProc();
        void setURL(const char* url){ cUrl = url;}
    };

    std::string _consumerKey;
    std::string _consumerSecret;
    std::string _accessTokenKey;
    std::string _accessTokenSecret;

    std::thread _runloop;
    std::atomic<bool> _stop{false};
    Proc _proc;
    std::string _url{"https://stream.twitter.com/1.1/statuses/filter.json?track="};
    LineFramer _framer{[this](std::string & docString){
        auto doc = BSON::Value::fromJSON(docString);
        emit(std::move(doc));
    }};
    void stop(){
        _stop.store(true);
    }
  public:  
    TwitterSource(std::shared_ptr<ProcessorManager> mgr, const std::string & id, 
        std::string filter,
        std::string consumerKey, std::string consumerSecret,
        std::string accessTokenKey, std::string accessTokenSecret) : 
            Processor{mgr,id},
            _consumerKey{consumerKey},
            _consumerSecret{consumerSecret},
            _accessTokenKey{accessTokenKey},
            _accessTokenSecret{accessTokenSecret},
            _proc{this,_consumerKey.c_str(),_consumerSecret.c_str(),accessTokenKey.c_str(),accessTokenSecret.c_str()} {
        _url += filter;
        _proc.setURL(_url.c_str());
        _runloop = std::move(std::thread{[this](){
            _proc.execProc();
        }});
    }

    //must be public to be accessable from curl -> do not use this.
    bool shouldStop(){
        return _stop.load();
    }

    void collect(char * data, std::size_t len){
        _framer.collect(data,len);
    }

    virtual ~TwitterSource(){
        stop();
        if(_runloop.joinable())_runloop.join();
    }
    
    virtual void process(BSON::Value && value, int input) override {}

};

size_t fncCallback(char* ptr, size_t size, size_t nmemb, TwitterSource *twitterSource) {
    size_t iRealSize = size * nmemb;
    twitterSource->collect(ptr,iRealSize);
    return iRealSize;
}

int progressCallback(TwitterSource *twitterSource, double dltotal, double dlnow, double ultotal, double ulnow){
    if(twitterSource->shouldStop()){
        return 1;
    }
    return 0;
}


// Constructor
TwitterSource::Proc::Proc( TwitterSource *twitterSource,
    const char* cConsKey, const char* cConsSec,
    const char* cAtokKey, const char* cAtokSec)
{
    this->_twitterSource = twitterSource;
    this->cConsKey = cConsKey;
    this->cConsSec = cConsSec;
    this->cAtokKey = cAtokKey;
    this->cAtokSec = cAtokSec;
}

void TwitterSource::Proc::execProc()
{
    // ==== cURL Initialization
    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (!curl) {
        std::cerr << "[ERROR] curl_easy_init" << std::endl;
        curl_global_cleanup();
        return;
    }

    // ==== cURL Setting
    // - URL, POST parameters, OAuth signing method, HTTP method, OAuth keys
    cSignedUrl = oauth_sign_url2(
        cUrl, NULL, OA_HMAC, "GET",
        cConsKey, cConsSec, cAtokKey, cAtokSec
    );
    // - URL
    curl_easy_setopt(curl, CURLOPT_URL, cSignedUrl);
    // - User agent name
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "mk-mode BOT");
    // - HTTP STATUS >=400 ---> ERROR
    curl_easy_setopt(curl, CURLOPT_FAILONERROR, 1);
    // - Callback function
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, fncCallback);
    // - Write data
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)_twitterSource);
    // - Progress callback (check if stop)
    curl_easy_setopt(curl, CURLOPT_PROGRESSFUNCTION, progressCallback);
    // - Write data
    curl_easy_setopt(curl, CURLOPT_PROGRESSDATA, (void *)_twitterSource);

    // ==== Execute
    int iStatus = curl_easy_perform(curl);
    if (!iStatus)
        std::cerr << "[ERROR] curl_easy_perform: STATUS=" << iStatus << std::endl;

    // ==== cURL Cleanup
    curl_easy_cleanup(curl);
    curl_global_cleanup();
}


}

#endif // __TWITTERSOURCE__