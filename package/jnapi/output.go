package jnapi

var CxxMain = `// region JN_ENTRY_POINT
int main(void) {
    std::set_terminate(&jn_terminate_handler);
    std::cout << std::boolalpha;
#ifdef _WINDOWS
    SetConsoleOutputCP(CP_UTF8);
    _setmode(_fileno(stdin), 0x00020000);
#endif
    ` + InitializerCaller + `();
    JNID(main());
    return EXIT_SUCCESS;
}
// endregion JN_ENTRY_POINT`

var CxxDefault = `#if defined(WIN32) || defined(_WIN32) || defined(__WIN32__) || defined(__NT__)
#define _WINDOWS
#endif

// region JN_STANDARD_IMPORTS
#include <iostream>
#include <cstring>
#include <string>
#include <sstream>
#include <functional>
#include <vector>
#include <map>
#include <thread>
#include <typeinfo>
#ifdef _WINDOWS
#include <codecvt>
#include <windows.h>
#include <fcntl.h>
#endif
// endregion JN_STANDARD_IMPORTS

#define _CONCAT(_A, _B) _A ## _B
#define CONCAT(_A, _B) _CONCAT(_A, _B)
#define JNID(_Identifier) CONCAT(_, _Identifier)

static inline void JNID(panic)(const char *_Message);

// region JN_CXX_API
// region JN_BUILTIN_VALUES
#define nil nullptr
// endregion JN_BUILTIN_VALUES

// region JN_BUILTIN_TYPES
typedef std::size_t                       uint_jnt;
typedef std::make_signed<uint_jnt>::type   int_jnt;
typedef int8_t                            i8_jnt;
typedef int16_t                           i16_jnt;
typedef int32_t                           i32_jnt;
typedef int64_t                           i64_jnt;
typedef uint8_t                           u8_jnt;
typedef uint16_t                          u16_jnt;
typedef uint32_t                          u32_jnt;
typedef uint64_t                          u64_jnt;
typedef float                             f32_jnt;
typedef double                            f64_jnt;
typedef unsigned char                     char_jnt;
typedef bool                              bool_jnt;
typedef void                              *voidptr_jnt;
typedef intptr_t                          intptr_jnt;
typedef uintptr_t                         uintptr_jnt;
#define func std::function

// region JN_STRUCTURES
template<typename _Item_t>
class array {
public:
    std::vector<_Item_t> _buffer{};

    array<_Item_t>(void) noexcept                       {}
    array<_Item_t>(const std::nullptr_t) noexcept       {}
    array<_Item_t>(const array<_Item_t>& _Src) noexcept { this->_buffer = _Src._buffer; }

    array<_Item_t>(const std::initializer_list<_Item_t> &_Src) noexcept
    { this->_buffer = std::vector<_Item_t>{_Src.begin(), _Src.end()}; }

    ~array<_Item_t>(void) noexcept { this->_buffer.clear(); }

    typedef _Item_t       *iterator;
    typedef const _Item_t *const_iterator;
    iterator begin(void) noexcept             { return &this->_buffer[0]; }
    const_iterator begin(void) const noexcept { return &this->_buffer[0]; }
    iterator end(void) noexcept               { return &this->_buffer[this->_buffer.size()]; }
    const_iterator end(void) const noexcept   { return &this->_buffer[this->_buffer.size()]; }

    inline void clear(void) noexcept        { this->_buffer.clear(); }
    inline uint_jnt len(void) const noexcept { return this->_buffer.size(); }
    inline bool empty(void) const noexcept  { return this->_buffer.empty(); }

    _Item_t *find(const _Item_t &_Item) noexcept {
        iterator _it{this->begin()};
        const iterator _end{this->end()};
        for (; _it < _end; ++_it)
        { if (_Item == *_it) { return _it; } }
        return nil;
    }

    _Item_t *rfind(const _Item_t &_Item) noexcept {
        iterator _it{this->end()};
        const iterator _begin{this->begin()};
        for (; _it >= _begin; --_it)
        { if (_Item == *_it) { return _it; } }
        return nil;
    }

    void erase(const _Item_t &_Item) noexcept {
        auto _it{this->_buffer.begin()};
        auto _end{this->_buffer.end()};
        for (; _it < _end; ++_it) {
            if (_Item == *_it) {
                this->_buffer.erase(_it);
                return;
            }
        }
    }

    void erase_all(const _Item_t &_Item) noexcept {
        auto _it{this->_buffer.begin()};
        auto _end{this->_buffer.end()};
        for (; _it < _end; ++_it)
        { if (_Item == *_it) { this->_buffer.erase(_it); } }
    }

    void append(const array<_Item_t> &_Items) noexcept
    { for (const _Item_t _item: _Items) { this->_buffer.push_back(_item); } }

    bool insert(const uint_jnt &_Start, const array<_Item_t> &_Items) noexcept {
        auto _it{this->_buffer.begin()+_Start};
        if (_it >= this->_buffer.end()) { return false; }
        this->_buffer.insert(_it, _Items.begin(), _Items.end());
        return true;
    }

    bool operator==(const array<_Item_t> &_Src) const noexcept {
        const uint_jnt _length{this->_buffer.size()};
        const uint_jnt _Src_length{_Src._buffer.size()};
        if (_length != _Src_length) { return false; }
        for (uint_jnt _index{0}; _index < _length; ++_index)
        { if (this->_buffer[_index] != _Src._buffer[_index]) { return false; } }
        return true;
    }

    bool operator!=(const array<_Item_t> &_Src) const noexcept { return !this->operator==(_Src); }
    bool operator==(const std::nullptr_t) const noexcept       { return this->_buffer.empty(); }
    bool operator!=(const std::nullptr_t) const noexcept       { return !this->operator==(nil); }

    _Item_t& operator[](const uint_jnt _Index) {
        if (this->len() <= _Index) { JNID(panic)("index out of range"); }
        return this->_buffer[_Index];
    }

    friend std::ostream& operator<<(std::ostream &_Stream,
                                    const array<_Item_t> &_Src) {
        _Stream << '[';
        const uint_jnt _length{_Src._buffer.size()};
        for (uint_jnt _index{0}; _index < _length;) {
            _Stream << _Src._buffer[_index++];
            if (_index < _length) { _Stream << u8", "; }
        }
        _Stream << ']';
        return _Stream;
    }
};
// endregion JN_STRUCTURES

template<typename _Key_t, typename _Value_t>
class map: public std::map<_Key_t, _Value_t> {
public:
    map<_Key_t, _Value_t>(void) noexcept                 {}
    map<_Key_t, _Value_t>(const std::nullptr_t) noexcept {}
    map<_Key_t, _Value_t>(const std::initializer_list<std::pair<_Key_t, _Value_t>> _Src)
    { for (const auto _data: _Src) { this->insert(_data); } }

    array<_Key_t> keys(void) const noexcept {
        array<_Key_t> _keys{};
        for (const auto _pair: *this)
        { _keys._buffer.push_back(_pair.first); }
        return _keys;
    }

    array<_Value_t> values(void) const noexcept {
        array<_Value_t> _values{};
        for (const auto _pair: *this)
        { _values._buffer.push_back(_pair.second); }
        return _values;
    }

    inline bool has(const _Key_t _Key) const noexcept { return this->find(_Key) != this->end(); }
    inline void del(const _Key_t _Key) noexcept { this->erase(_Key); }

    bool operator==(const std::nullptr_t) const noexcept { return this->empty(); }
    bool operator!=(const std::nullptr_t) const noexcept { return !this->operator==(nil); }

    friend std::ostream& operator<<(std::ostream &_Stream,
                                    const map<_Key_t, _Value_t> &_Src) {
        _Stream << '{';
        uint_jnt _length{_Src.size()};
        for (const auto _pair: _Src) {
            _Stream << _pair.first;
            _Stream << ':';
            _Stream << _pair.second;
            if (--_length > 0) { _Stream << u8", "; }
        }
        _Stream << '}';
        return _Stream;
    }
};

class str_jnt {
public:
    std::string _buffer{};

    str_jnt(void) noexcept                   {}
    str_jnt(const char *_Src) noexcept       { this->_buffer = _Src ? _Src : ""; }
    str_jnt(const std::string _Src) noexcept { this->_buffer = _Src; }
    str_jnt(const str_jnt &_Src) noexcept     { this->_buffer = _Src._buffer; }

    str_jnt(const array<char> &_Src) noexcept
    { this->_buffer = std::string{_Src.begin(), _Src.end()}; }

    str_jnt(const array<u8_jnt> &_Src) noexcept
    { this->_buffer = std::string{_Src.begin(), _Src.end()}; }

    typedef char_jnt       *iterator;
    typedef const char_jnt *const_iterator;
    iterator begin(void) noexcept             { return (iterator)(&this->_buffer[0]); }
    const_iterator begin(void) const noexcept { return (const_iterator)(&this->_buffer[0]); }
    iterator end(void) noexcept               { return (iterator)(&this->_buffer[this->len()]); }
    const_iterator end(void) const noexcept   { return (const_iterator)(&this->_buffer[this->len()]); }

    inline uint_jnt len(void) const noexcept { return this->_buffer.length(); }
    inline bool empty(void) const noexcept  { return this->_buffer.empty(); }

    inline str_jnt sub(const uint_jnt start, const uint_jnt end) const noexcept
    { return this->_buffer.substr(start, end); }

    inline str_jnt sub(const uint_jnt start) const noexcept
    { return this->_buffer.substr(start); }

    inline bool has_prefix(const str_jnt &_Sub) const noexcept
    { return this->len() >= _Sub.len() && this->sub(0, _Sub.len()) == _Sub._buffer; }

    inline bool has_suffix(const str_jnt &_Sub) const noexcept
    { return this->len() >= _Sub.len() && this->sub(this->len()-_Sub.len()) == _Sub; }

    inline uint_jnt find(const str_jnt &_Sub) const noexcept
    { return this->_buffer.find(_Sub._buffer); }

    inline uint_jnt rfind(const str_jnt &_Sub) const noexcept
    { return this->_buffer.rfind(_Sub._buffer); }

    inline const char* cstr(void) const noexcept
    { return this->_buffer.c_str(); }

    str_jnt trim(const str_jnt &_Bytes) const noexcept {
        const_iterator _it{this->begin()};
        const const_iterator _end{this->end()};
        const_iterator _begin{this->begin()};
        for (; _it < _end; ++_it) {
            bool exist{false};
            const_iterator _bytes_it{_Bytes.begin()};
            const const_iterator _bytes_end{_Bytes.end()};
            for (; _bytes_it < _bytes_end; ++_bytes_it)
            { if ((exist = *_it == *_bytes_it)) { break; } }
            if (!exist) { return this->sub(_it-_begin); }
        }
        return str_jnt{u8""};
    }

    str_jnt rtrim(const str_jnt &_Bytes) const noexcept {
        const_iterator _it{this->end()-1};
        const const_iterator _begin{this->begin()};
        for (; _it >= _begin; --_it) {
            bool exist{false};
            const_iterator _bytes_it{_Bytes.begin()};
            const const_iterator _bytes_end{_Bytes.end()};
            for (; _bytes_it < _bytes_end; ++_bytes_it)
            { if ((exist = *_it == *_bytes_it)) { break; } }
            if (!exist) { return this->sub(0, _it-_begin+1); }
        }
        return str_jnt{u8""};
    }

    array<str_jnt> split(const str_jnt &_Sub, const i64_jnt &_N) const noexcept {
        array<str_jnt> _parts{};
        if (_N == 0) { return _parts; }
        const const_iterator _begin{this->begin()};
        std::string _s{this->_buffer};
        uint_jnt _pos{std::string::npos};
        if (_N < 0) {
            while ((_pos = _s.find(_Sub._buffer)) != std::string::npos) {
                _parts._buffer.push_back(_s.substr(0, _pos));
                _s = _s.substr(_pos+_Sub.len());
            }
            if (!_parts.empty()) { _parts._buffer.push_back(str_jnt{_s}); }
        } else {
            uint_jnt _n{0};
            while ((_pos = _s.find(_Sub._buffer)) != std::string::npos) {
                _parts._buffer.push_back(_s.substr(0, _pos));
                _s = _s.substr(_pos+_Sub.len());
                if (++_n >= _N) { break; }
            }
            if (!_parts.empty() && _n < _N) { _parts._buffer.push_back(str_jnt{_s}); }
        }
        return _parts;
    }

    str_jnt replace(const str_jnt &_Sub,
                   const str_jnt &_New,
                   const i64_jnt &_N) const noexcept {
        if (_N == 0) { return *this; }
        std::string _s{this->_buffer};
        uint_jnt start_pos{0};
        if (_N < 0) {
            while((start_pos = _s.find(_Sub._buffer, start_pos)) != std::string::npos) {
                _s.replace(start_pos, _Sub.len(), _New._buffer);
                start_pos += _New.len();
            }
        } else {
            uint_jnt _n{0};
            while((start_pos = _s.find(_Sub._buffer, start_pos)) != std::string::npos) {
                _s.replace(start_pos, _Sub.len(), _New._buffer);
                start_pos += _New.len();
                if (++_n >= _N) { break; }
            }
        }
        return str_jnt{_s};
    }

    operator array<char>(void) const noexcept {
        array<char> _array{};
        _array._buffer = std::vector<char>{this->begin(), this->end()};
        return _array;
    }

    operator array<u8_jnt>(void) const noexcept {
        array<u8_jnt> _array{};
        _array._buffer = std::vector<u8_jnt>{this->begin(), this->end()};
        return _array;
    }

    operator const char*(void) const noexcept
    { return this->_buffer.c_str(); }

    operator char*(void) const noexcept
    { return (char*)(this->_buffer.c_str()); }

    char &operator[](uint_jnt _Index) {
        if (this->len() <= _Index) { JNID(panic)("index out of range"); }
        return this->_buffer[_Index];
    }

    void operator+=(const str_jnt _Str) noexcept        { this->_buffer += _Str._buffer; }
    str_jnt operator+(const str_jnt _Str) const noexcept { return str_jnt{this->_buffer + _Str._buffer}; }
    bool operator==(const str_jnt &_Str) const noexcept { return this->_buffer == _Str._buffer; }
    bool operator!=(const str_jnt &_Str) const noexcept { return !this->operator==(_Str); }

    friend std::ostream& operator<<(std::ostream &_Stream, const str_jnt &_Src)
    { return _Stream << _Src._buffer; }
};

struct any_jnt {
public:
    void *_expr{nil};
    char *_inf{nil};

    template<typename T>
    any_jnt(const T &_Expr) noexcept
    { this->operator=(_Expr); }

    ~any_jnt(void) noexcept
    { this->_delete(); }

    inline void _delete(void) noexcept {
        this->_expr = nil;
        this->_inf = nil;
    }

    template<typename T>
    inline bool type_is(void) const noexcept {
        if (!this->_expr)
        { return std::is_same<std::nullptr_t, T>::value; }
        return std::strcmp(this->_inf, typeid(T).name()) == 0;
    }

    template<typename T>
    void operator=(const T &_Expr) noexcept {
        this->_delete();
        this->_expr = (void*)&_Expr;
        this->_inf  = (char*)(typeid(T).name());
    }

    void operator=(const std::nullptr_t) noexcept
    { this->_delete(); }

    template<typename T>
    operator T(void) const noexcept {
        if (!this->_expr)
        { JNID(panic)("casting failed because data is nil"); }
        if (std::strcmp(this->_inf, typeid(T).name()) != 0)
        { JNID(panic)("incompatible type"); }
        return *(T*)(this->_expr);
    }

    template<typename T>
    inline bool operator==(const T &_Expr) const noexcept
    { return this->type_is<T>() && *(T*)(this->_expr) == _Expr; }

    template<typename T>
    inline bool operator!=(const T &_Expr) const noexcept
    { return !this->operator==(_Expr); }

    inline bool operator==(const any_jnt &_Any) const noexcept
    { return this->_expr == _Any._expr; }

    inline bool operator!=(const any_jnt &_Any) const noexcept
    { return !this->operator==(_Any); }

    friend std::ostream& operator<<(std::ostream &_Stream, const any_jnt &_Src)
    { return _Stream << _Src._expr; }
};
// endregion JN_BUILTIN_TYPES

// region JN_MISC
template <typename _Enum_t, typename _Index_t, typename _Item_t>
static inline void foreach(const _Enum_t _Enum,
                           const func<void(_Index_t, _Item_t)> _Body) {
    _Index_t _index{0};
    for (auto _item: _Enum) { _Body(_index++, _item); }
}

template <typename _Enum_t, typename _Index_t>
static inline void foreach(const _Enum_t _Enum,
                           const func<void(_Index_t)> _Body) {
    _Index_t _index{0};
    for (auto begin = _Enum.begin(), end = _Enum.end(); begin < end; ++begin)
    { _Body(_index++); }
}

template <typename _Key_t, typename _Value_t>
static inline void foreach(const map<_Key_t, _Value_t> _Map,
                           const func<void(_Key_t)> _Body) {
    for (const auto _pair: _Map) { _Body(_pair.first); }
}

template <typename _Key_t, typename _Value_t>
static inline void foreach(const map<_Key_t, _Value_t> _Map,
                           const func<void(_Key_t, _Value_t)> _Body) {
    for (const auto _pair: _Map) { _Body(_pair.first, _pair.second); }
}

template <typename ...T>
static inline std::string strpol(const T... _Expression) noexcept {
  return (std::stringstream{} << ... << _Expression).str();
}

template<typename Type, unsigned N, unsigned Last>
struct tuple_ostream {
    static void arrow(std::ostream &_Stream, const Type &_Type) {
        _Stream << std::get<N>(_Type) << u8", ";
        tuple_ostream<Type, N + 1, Last>::arrow(_Stream, _Type);
    }
};

template<typename Type, unsigned N>
struct tuple_ostream<Type, N, N> {
    static void arrow(std::ostream &_Stream, const Type &_Type)
    { _Stream << std::get<N>(_Type); }
};

template<typename... Types>
std::ostream& operator<<(std::ostream &_Stream,
                         const std::tuple<Types...> &_Tuple) {
    _Stream << u8"(";
    tuple_ostream<std::tuple<Types...>, 0, sizeof...(Types)-1>::arrow(_Stream, _Tuple);
    _Stream << u8")";
    return _Stream;
}

template<typename _Function_t, typename _Tuple_t, size_t ... _I_t>
inline auto tuple_as_args(const _Function_t _Function,
                          const _Tuple_t _Tuple,
                          const std::index_sequence<_I_t ...>)
{ return _Function(std::get<_I_t>(_Tuple) ...); }

template<typename _Function_t, typename _Tuple_t>
inline auto tuple_as_args(const _Function_t _Function, const _Tuple_t _Tuple) {
    static constexpr auto _size{std::tuple_size<_Tuple_t>::value};
    return tuple_as_args(_Function, _Tuple, std::make_index_sequence<_size>{});
}

struct defer {
    typedef func<void(void)> _Function_t;
    template<class Callable>
    defer(Callable &&_function): _function(std::forward<Callable>(_function)) {}
    defer(defer &&_Src): _function(std::move(_Src._function))                 { _Src._function = nullptr; }
    ~defer() noexcept                                                         { if (this->_function) { this->_function(); } }
    defer(const defer &)          = delete;
    void operator=(const defer &) = delete;
    _Function_t _function;
};

std::ostream &operator<<(std::ostream &_Stream, const i8_jnt &_Src)
{ return _Stream << (i32_jnt)(_Src); }

std::ostream &operator<<(std::ostream &_Stream, const u8_jnt &_Src)
{ return _Stream << (i32_jnt)(_Src); }

std::ostream &operator<<(std::ostream &_Stream, const std::nullptr_t)
{ return _Stream << "<nil>"; }

template<typename _Obj_t>
str_jnt tostr(const _Obj_t &_Obj) noexcept {
    std::stringstream _stream;
    _stream << _Obj;
    return str_jnt{_stream.str()};
}

#define DEFER(_Expr) defer CONCAT(JNXDEFER_, __LINE__){[&](void) mutable -> void { _Expr; }}
#define CO(_Expr) std::thread{[&](void) mutable -> void { _Expr; }}.detach()
// endregion JN_MISC

// region PANIC_DEFINES
struct JNID(Error) {
public:
  str_jnt JNID(message);
};

std::ostream &operator<<(std::ostream &_Stream, const JNID(Error) &_Error)
{ return _Stream << _Error.JNID(message); }

static inline void JNID(panic)(const struct JNID(Error) &_Error) { throw _Error; }
static inline void JNID(panic)(const char *_Message) { JNID(panic)(JNID(Error){_Message}); }
// endregion PANIC_DEFINES

// region JN_BUILTIN_FUNCTIONS
template<typename _Obj_t>
static inline void JNID(print)(const _Obj_t _Obj) noexcept { std::cout << _Obj; }

template<typename _Obj_t>
static inline void JNID(println)(const _Obj_t _Obj) noexcept {
    JNID(print)<_Obj_t>(_Obj);
    std::cout << std::endl;
}
// endregion JN_BUILTIN_FUNCTIONS

// region BOTTOM_MISC
void jn_terminate_handler(void) noexcept {
    try { std::rethrow_exception(std::current_exception()); }
    catch (const JNID(Error) _error)
    { std::cout << "panic: " << _error.JNID(message) << std::endl; }
    catch (...)
    { std::cout << "panic: <undefined panics>" << std::endl; }
    std::exit(EXIT_FAILURE);
}
// endregion BOTTOM_MISC
`
