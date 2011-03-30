MethodDecl:
   void methodName(String[] args){
   }

expect:
    METHOD(
        MODIFIERS,
        TYPE('void'),
        IDENT('methodName'),
        ARGS(
            ARG(TYPE('String',DIM('1')),IDENT('args'),<nil>)
        ),
        METHOD_BODY
    )
