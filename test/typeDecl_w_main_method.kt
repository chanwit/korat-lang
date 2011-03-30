TypeDecl:

    class A {
       static main(args){
       }
    }
    
expect:
    CLASS(
        IDENT('A'),
        MEMBERS(
            METHOD(
                MODIFIERS(STATIC),
                <nil>,
                IDENT('main'),
                ARGS(
                    ARG(TYPE('java.lang.Object'),IDENT('args'),<nil>)
                ),
                METHOD_BODY
            )
        )
    )
