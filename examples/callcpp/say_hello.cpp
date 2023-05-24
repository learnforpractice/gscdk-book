#include <eosio/print.hpp>

extern "C" void say_hello(const char *s, uint32_t len) {
    std::string _s(s, len);
    eosio::print(_s);
}
