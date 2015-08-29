cmake_minimum_required (VERSION 2.8.11)

project (horst-twittersource)

set( CMAKE_COLOR_MAKEFILE ON )
include_directories( ${CMAKE_CURRENT_SOURCE_DIR}/include )

set( SRC ${CMAKE_CURRENT_SOURCE_DIR}/source/export.cpp
		${CMAKE_CURRENT_SOURCE_DIR}/source/LineFramer.cpp )

set( CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -fPIC --std=c++11 -Os -O3" )

add_library( ${PROJECT_NAME} SHARED ${SRC} )
target_link_libraries( ${PROJECT_NAME} curl oauth pthread)