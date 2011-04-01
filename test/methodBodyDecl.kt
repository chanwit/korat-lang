MethodBodyDecl:
    {
        a := match {
        }
    }

expect:
    METHOD_BODY(
        INFER_ASSIGN(
            LOCAL_VAR('a'),
            MATCH
        )
    )
