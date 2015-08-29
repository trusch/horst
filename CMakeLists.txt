cmake_minimum_required (VERSION 2.8.11)

project (horst)

set( CMAKE_BUILD_TYPE Debug )

set( CMAKE_COLOR_MAKEFILE ON )
set( CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin/plugins )
set( CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin )

add_subdirectory(${CMAKE_CURRENT_SOURCE_DIR}/runner)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/runner/include)
add_subdirectory(${CMAKE_CURRENT_SOURCE_DIR}/plugins/utils)
add_subdirectory(${CMAKE_CURRENT_SOURCE_DIR}/plugins/twittersource)