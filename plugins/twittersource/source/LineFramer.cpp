/*
 * Copyright (c) 2015, webvariants GmbH, http://www.webvariants.de
 *
 * This file is released under the terms of the MIT license. You can find the
 * complete text in the attached LICENSE file or online at:
 *
 * http://www.opensource.org/licenses/mit-license.php
 *
 * @author: Tino Rusch (tino.rusch@webvariants.de)
 */

#include "horst/plugins/LineFramer.h" 

Horst::LineFramer::LineFramer( std::function<void( std::string& )> onMessage ) : _onMessage {onMessage} {}
Horst::LineFramer::LineFramer() {}

void Horst::LineFramer::operator=( const LineFramer & other ) {
    _onMessage = other._onMessage;
};

void Horst::LineFramer::collect( char * data, size_t len ) {
    for( size_t i=0; i<len; i++ ) {
        char c = data[i];
        message += c;
        if(c=='\n' && message.size()>=2 && message[message.size()-2] == '\r') {
            _onMessage(message);
            message = "";
        }
    }
}
