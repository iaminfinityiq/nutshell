int a = 5
double b = 6.0
string c = "These are one of the most commonly used data types in this language"
void d = null

type e = int
any f = "e is a type object, and i am a string object, but since every data type inherits the any type, then this is still legal"
builtin_function g = println

g("int:", a)
g("double:", b)
g("string:", c)
g("void:", d)

g("type:", e)
g("any:", f)
g("builtin_function:", g)