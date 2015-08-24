cmake_minimum_required (VERSION 2.8.11)

project (horst-runner)

include_directories( ${CMAKE_CURRENT_SOURCE_DIR}/include )

set( SRC ${CMAKE_CURRENT_SOURCE_DIR}/source/ProcessorManager.cpp
			${CMAKE_CURRENT_SOURCE_DIR}/source/main.cpp )

set( CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall --std=c++11 -Os -O3" )

add_executable( ${PROJECT_NAME} ${SRC} )
target_link_libraries( ${PROJECT_NAME} pthread )
target_link_libraries( ${PROJECT_NAME} bson )
target_link_libraries( ${PROJECT_NAME} dl )
